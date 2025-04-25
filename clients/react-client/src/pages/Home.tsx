import { PackageTracker } from '../components/PackageTracker';

export function Home() {
  return (
    <div className="space-y-6">
      <div className="bg-white shadow rounded-lg p-6">
        <h2 className="text-xl font-semibold text-gray-900 mb-4">Welcome to MicroTracker</h2>
        <p className="text-gray-600">
          Track and monitor your microservices with ease. Get real-time insights into your system's performance and health.
        </p>
      </div>
      
      <PackageTracker />
      
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div className="bg-white shadow rounded-lg p-6">
          <h3 className="text-lg font-medium text-gray-900 mb-2">Service Health</h3>
          <p className="text-gray-600">Monitor the health status of all your microservices in one place.</p>
        </div>
        
        <div className="bg-white shadow rounded-lg p-6">
          <h3 className="text-lg font-medium text-gray-900 mb-2">Performance Metrics</h3>
          <p className="text-gray-600">Track response times, error rates, and other key performance indicators.</p>
        </div>
        
        <div className="bg-white shadow rounded-lg p-6">
          <h3 className="text-lg font-medium text-gray-900 mb-2">Logs & Traces</h3>
          <p className="text-gray-600">View and analyze logs and traces across your microservices architecture.</p>
        </div>
      </div>
    </div>
  );
} 