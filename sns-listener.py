#!/usr/bin/env python3

from flask import Flask, request, jsonify, render_template_string
import json
import logging
import datetime
import sys
import boto3
import threading
import time

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    stream=sys.stdout
)
logger = logging.getLogger('sns-listener')

app = Flask(__name__)

# In-memory storage for messages
messages = []
topics = {}  # To store topic details

# HTML template for the web interface
HTML_TEMPLATE = """
<!DOCTYPE html>
<html>
<head>
    <title>SNS Message Listener</title>
    <meta http-equiv="refresh" content="5">
    <style>
        body {
            font-family: Arial, sans-serif;
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px;
        }
        h1, h2 {
            color: #232f3e;
        }
        .message {
            border: 1px solid #ddd;
            padding: 15px;
            margin: 10px 0;
            border-radius: 4px;
            background-color: #f9f9f9;
        }
        .message pre {
            white-space: pre-wrap;
            word-wrap: break-word;
        }
        .timestamp {
            color: #666;
            font-size: 0.8em;
        }
        .topic {
            font-weight: bold;
            color: #c15317;
        }
        .clear-button {
            background-color: #ff9900;
            color: white;
            border: none;
            padding: 10px 15px;
            cursor: pointer;
            margin: 5px 0;
            border-radius: 3px;
        }
        .stats {
            margin-bottom: 20px;
            padding: 10px;
            background-color: #eee;
            border-radius: 4px;
        }
    </style>
</head>
<body>
    <h1>SNS Message Listener</h1>
    
    <div class="stats">
        <h3>Statistics</h3>
        <p>Total messages received: {{ messages|length }}</p>
        <p>Topics: {{ topics|length }}</p>
        <ul>
            {% for topic_name, count in topics.items() %}
            <li>{{ topic_name }}: {{ count }} message(s)</li>
            {% endfor %}
        </ul>
    </div>
    
    <h2>Messages</h2>
    <form method="POST" action="/clear">
        <button type="submit" class="clear-button">Clear Messages</button>
    </form>
    
    {% if messages %}
        {% for msg in messages|reverse %}
        <div class="message">
            <p class="timestamp">{{ msg.timestamp }}</p>
            <p class="topic">Topic: {{ msg.topic_name }}</p>
            <pre>{{ msg.content }}</pre>
        </div>
        {% endfor %}
    {% else %}
        <p>No messages received yet. Publish a message to your SNS topic!</p>
    {% endif %}
</body>
</html>
"""

@app.route('/')
def home():
    """Display the web interface with received messages."""
    return render_template_string(HTML_TEMPLATE, messages=messages, topics=topics)

@app.route('/clear', methods=['POST'])
def clear_messages():
    """Clear all stored messages."""
    global messages, topics
    messages = []
    topics = {}
    return render_template_string(HTML_TEMPLATE, messages=messages, topics=topics)

@app.route('/webhook', methods=['POST', 'GET'])
def webhook():
    """Receive SNS notifications."""
    # Handle SNS subscription confirmation
    if request.method == 'GET':
        return "SNS Webhook is active"
    
    try:
        data = request.json
        logger.info(f"Received notification: {json.dumps(data, indent=2)}")
        
        # Handle SNS subscription confirmation
        if 'Type' in data and data['Type'] == 'SubscriptionConfirmation':
            confirm_subscription(data)
            return jsonify({"status": "subscription_confirming"})
        
        # Handle SNS message
        if 'Type' in data and data['Type'] == 'Notification':
            topic_arn = data.get('TopicArn', 'Unknown')
            topic_name = topic_arn.split(':')[-1]
            
            # Update topic stats
            if topic_name in topics:
                topics[topic_name] += 1
            else:
                topics[topic_name] = 1
            
            # Store message
            message_data = {
                'timestamp': datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S'),
                'topic_arn': topic_arn,
                'topic_name': topic_name,
                'message_id': data.get('MessageId', 'Unknown'),
                'content': data.get('Message', 'No content')
            }
            
            messages.append(message_data)
            logger.info(f"✅ Stored message from topic '{topic_name}': {message_data['content']}")
            
            return jsonify({"status": "message_received"})
    
    except Exception as e:
        logger.error(f"Error processing request: {e}")
        return jsonify({"status": "error", "message": str(e)}), 500
    
    return jsonify({"status": "unknown_request_type"})

def confirm_subscription(data):
    """Confirm SNS subscription by visiting the SubscribeURL."""
    import urllib.request
    
    subscribe_url = data.get('SubscribeURL')
    if subscribe_url:
        logger.info(f"Confirming subscription: {subscribe_url}")
        try:
            urllib.request.urlopen(subscribe_url)
            logger.info("✅ Subscription confirmed!")
        except Exception as e:
            logger.error(f"Failed to confirm subscription: {e}")

def subscribe_to_topic(topic_arn, endpoint):
    """Subscribe the endpoint to an SNS topic."""
    try:
        sns = boto3.client(
            'sns',
            endpoint_url='http://localhost:4566',
            region_name='us-east-1',
            aws_access_key_id='test',
            aws_secret_access_key='test'
        )
        
        response = sns.subscribe(
            TopicArn=topic_arn,
            Protocol='http',
            Endpoint=endpoint
        )
        logger.info(f"✅ Subscribed to topic {topic_arn}: {response}")
        return response
    except Exception as e:
        logger.error(f"Failed to subscribe to topic: {e}")
        return None

def list_topics():
    """List all SNS topics."""
    try:
        sns = boto3.client(
            'sns',
            endpoint_url='http://localhost:4566',
            region_name='us-east-1',
            aws_access_key_id='test',
            aws_secret_access_key='test'
        )
        
        response = sns.list_topics()
        return response.get('Topics', [])
    except Exception as e:
        logger.error(f"Failed to list topics: {e}")
        return []

def auto_subscribe(endpoint_url):
    """Automatically subscribe to all topics."""
    topics = list_topics()
    for topic in topics:
        topic_arn = topic.get('TopicArn')
        logger.info(f"Auto-subscribing to topic: {topic_arn}")
        subscribe_to_topic(topic_arn, endpoint_url)

if __name__ == '__main__':
    # Get all arguments
    import argparse
    parser = argparse.ArgumentParser(description='SNS Listener')
    parser.add_argument('--port', type=int, default=8000, help='Port to run the server on')
    parser.add_argument('--host', type=str, default='0.0.0.0', help='Host to bind to')
    parser.add_argument('--public-url', type=str, help='Public URL for subscription (e.g., http://example.com:8000/webhook)')
    args = parser.parse_args()
    
    port = args.port
    host = args.host
    
    # Determine public endpoint for subscription
    if args.public_url:
        public_endpoint = args.public_url
    else:
        public_endpoint = f"http://localhost:{port}/webhook"
    
    # Auto-subscribe in a separate thread to avoid blocking startup
    def subscription_worker():
        logger.info("Waiting for server to start before subscribing...")
        time.sleep(2)  # Give Flask time to start
        logger.info(f"Auto-subscribing to all topics using endpoint: {public_endpoint}")
        auto_subscribe(public_endpoint)
    
    threading.Thread(target=subscription_worker).start()
    
    # Print usage instructions
    logger.info("=" * 70)
    logger.info(f"SNS Listener starting on http://{host}:{port}")
    logger.info(f"Web interface: http://localhost:{port}")
    logger.info(f"Webhook endpoint: http://localhost:{port}/webhook")
    logger.info("=" * 70)
    
    # Start Flask server
    app.run(host=host, port=port, debug=True, use_reloader=False)

