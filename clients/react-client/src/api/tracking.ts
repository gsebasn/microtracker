import { TrackingResponse } from '../types/tracking';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

export const trackingApi = {
  getTrackingInfo: async (packageId: string): Promise<TrackingResponse> => {
    const response = await fetch(`${API_BASE_URL}/api/v1/tracking/${packageId}`);
    if (!response.ok) {
      throw new Error('Failed to fetch tracking information');
    }
    return response.json();
  },

  getTrackingHistory: async (packageId: string): Promise<TrackingResponse> => {
    const response = await fetch(`${API_BASE_URL}/api/v1/tracking/${packageId}/history`);
    if (!response.ok) {
      throw new Error('Failed to fetch tracking history');
    }
    return response.json();
  },
}; 