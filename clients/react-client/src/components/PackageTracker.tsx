import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { packageService } from '../api/services/packageService';

interface PackageTrackerProps {
  initialTrackingNumber?: string;
}

export function PackageTracker({ initialTrackingNumber }: PackageTrackerProps) {
  const [trackingNumber, setTrackingNumber] = useState(initialTrackingNumber || '');
  const [searchQuery, setSearchQuery] = useState('');

  const { data: packageData, isLoading: isLoadingPackage } = useQuery({
    queryKey: ['package', trackingNumber],
    queryFn: () => packageService.getPackage(trackingNumber),
    enabled: !!trackingNumber,
  });

  const { data: searchResults, isLoading: isLoadingSearch } = useQuery({
    queryKey: ['search', searchQuery],
    queryFn: () => packageService.searchPackages(searchQuery),
    enabled: !!searchQuery,
  });

  return (
    <div className="space-y-6">
      <div className="flex gap-4">
        <div className="flex-1">
          <label htmlFor="tracking" className="block text-sm font-medium text-gray-700">
            Tracking Number
          </label>
          <div className="mt-1">
            <input
              type="text"
              name="tracking"
              id="tracking"
              className="block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
              value={trackingNumber}
              onChange={(e) => setTrackingNumber(e.target.value)}
              placeholder="Enter tracking number"
            />
          </div>
        </div>
        <div className="flex-1">
          <label htmlFor="search" className="block text-sm font-medium text-gray-700">
            Search Packages
          </label>
          <div className="mt-1">
            <input
              type="text"
              name="search"
              id="search"
              className="block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              placeholder="Search packages..."
            />
          </div>
        </div>
      </div>

      {isLoadingPackage && (
        <div className="text-center py-4">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-500 mx-auto"></div>
          <p className="mt-2 text-sm text-gray-500">Loading package information...</p>
        </div>
      )}

      {packageData && (
        <div className="bg-white shadow rounded-lg p-6">
          <h3 className="text-lg font-medium text-gray-900 mb-4">Package Information</h3>
          <dl className="grid grid-cols-1 gap-x-4 gap-y-6 sm:grid-cols-2">
            <div>
              <dt className="text-sm font-medium text-gray-500">Tracking Number</dt>
              <dd className="mt-1 text-sm text-gray-900">{packageData.trackingNumber}</dd>
            </div>
            <div>
              <dt className="text-sm font-medium text-gray-500">Status</dt>
              <dd className="mt-1 text-sm text-gray-900">{packageData.status}</dd>
            </div>
            <div>
              <dt className="text-sm font-medium text-gray-500">Last Update</dt>
              <dd className="mt-1 text-sm text-gray-900">{packageData.lastUpdate}</dd>
            </div>
            {packageData.estimatedDelivery && (
              <div>
                <dt className="text-sm font-medium text-gray-500">Estimated Delivery</dt>
                <dd className="mt-1 text-sm text-gray-900">{packageData.estimatedDelivery}</dd>
              </div>
            )}
            {packageData.currentLocation && (
              <div>
                <dt className="text-sm font-medium text-gray-500">Current Location</dt>
                <dd className="mt-1 text-sm text-gray-900">{packageData.currentLocation}</dd>
              </div>
            )}
          </dl>
        </div>
      )}

      {isLoadingSearch && (
        <div className="text-center py-4">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-500 mx-auto"></div>
          <p className="mt-2 text-sm text-gray-500">Searching packages...</p>
        </div>
      )}

      {searchResults && searchResults.length > 0 && (
        <div className="bg-white shadow rounded-lg p-6">
          <h3 className="text-lg font-medium text-gray-900 mb-4">Search Results</h3>
          <div className="space-y-4">
            {searchResults.map((pkg) => (
              <div key={pkg.id} className="border-b border-gray-200 pb-4 last:border-0 last:pb-0">
                <div className="flex justify-between">
                  <div>
                    <p className="text-sm font-medium text-gray-900">{pkg.trackingNumber}</p>
                    <p className="text-sm text-gray-500">{pkg.status}</p>
                  </div>
                  <button
                    onClick={() => setTrackingNumber(pkg.trackingNumber)}
                    className="text-sm text-indigo-600 hover:text-indigo-900"
                  >
                    View Details
                  </button>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
} 