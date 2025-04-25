export interface PackageStatus {
  package_id: string;
  status: string;
  location: string;
  timestamp: string;
  description: string;
}

export interface TrackingInfo {
  package_id: string;
  current_status: PackageStatus;
  history: PackageStatus[];
  estimated_delivery?: string;
  carrier?: string;
  tracking_number?: string;
}

export interface TrackingResponse {
  success: boolean;
  data: TrackingInfo;
  error?: string;
} 