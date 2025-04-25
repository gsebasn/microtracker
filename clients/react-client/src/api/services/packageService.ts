import { apiClient } from '../client';

export interface Package {
  id: string;
  trackingNumber: string;
  status: string;
  lastUpdate: string;
  estimatedDelivery?: string;
  currentLocation?: string;
}

export const packageService = {
  async getPackage(trackingNumber: string): Promise<Package> {
    const response = await apiClient.get<Package>(`/packages/${trackingNumber}`);
    return response.data;
  },

  async getPackageHistory(trackingNumber: string): Promise<Package[]> {
    const response = await apiClient.get<Package[]>(`/packages/${trackingNumber}/history`);
    return response.data;
  },

  async searchPackages(query: string): Promise<Package[]> {
    const response = await apiClient.get<Package[]>('/packages/search', {
      params: { q: query }
    });
    return response.data;
  }
}; 