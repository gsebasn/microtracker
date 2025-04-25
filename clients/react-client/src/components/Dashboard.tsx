import { useQuery } from '@tanstack/react-query';
import { trackingApi } from '../api/tracking';
import { format } from 'date-fns';
import { PackageStatus } from '../types/tracking';

const StatusBadge = ({ status }: { status: string }) => {
  const getStatusColor = (status: string) => {
    switch (status.toLowerCase()) {
      case 'delivered':
        return 'bg-green-100 text-green-800';
      case 'in_transit':
        return 'bg-blue-100 text-blue-800';
      case 'pending':
        return 'bg-yellow-100 text-yellow-800';
      case 'failed':
        return 'bg-red-100 text-red-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  return (
    <span className={`px-2 py-1 rounded-full text-sm font-medium ${getStatusColor(status)}`}>
      {status.replace('_', ' ')}
    </span>
  );
};

const StatusTimeline = ({ history }: { history: PackageStatus[] }) => {
  return (
    <div className="flow-root">
      <ul role="list" className="-mb-8">
        {history.map((status, index) => (
          <li key={index}>
            <div className="relative pb-8">
              {index !== history.length - 1 && (
                <span
                  className="absolute left-4 top-4 -ml-px h-full w-0.5 bg-gray-200"
                  aria-hidden="true"
                />
              )}
              <div className="relative flex space-x-3">
                <div>
                  <span className="h-8 w-8 rounded-full bg-gray-400 flex items-center justify-center ring-8 ring-white">
                    <svg
                      className="h-5 w-5 text-white"
                      viewBox="0 0 20 20"
                      fill="currentColor"
                      aria-hidden="true"
                    >
                      <path
                        fillRule="evenodd"
                        d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z"
                        clipRule="evenodd"
                      />
                    </svg>
                  </span>
                </div>
                <div className="flex min-w-0 flex-1 justify-between space-x-4 pt-1.5">
                  <div>
                    <p className="text-sm text-gray-500">
                      {status.description}{' '}
                      <span className="font-medium text-gray-900">{status.location}</span>
                    </p>
                  </div>
                  <div className="whitespace-nowrap text-right text-sm text-gray-500">
                    <time dateTime={status.timestamp}>
                      {format(new Date(status.timestamp), 'MMM d, yyyy h:mm a')}
                    </time>
                  </div>
                </div>
              </div>
            </div>
          </li>
        ))}
      </ul>
    </div>
  );
};

export const Dashboard = ({ packageId }: { packageId: string }) => {
  const { data, isLoading, error } = useQuery({
    queryKey: ['tracking', packageId],
    queryFn: () => trackingApi.getTrackingInfo(packageId),
  });

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-gray-900"></div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="rounded-md bg-red-50 p-4">
        <div className="flex">
          <div className="ml-3">
            <h3 className="text-sm font-medium text-red-800">Error loading tracking information</h3>
            <div className="mt-2 text-sm text-red-700">
              <p>{(error as Error).message}</p>
            </div>
          </div>
        </div>
      </div>
    );
  }

  if (!data?.data) {
    return null;
  }

  const { current_status, history, estimated_delivery, carrier, tracking_number } = data.data;

  return (
    <div className="bg-white shadow sm:rounded-lg">
      <div className="px-4 py-5 sm:p-6">
        <div className="grid grid-cols-1 gap-6 sm:grid-cols-2">
          <div>
            <h3 className="text-lg font-medium leading-6 text-gray-900">Package Information</h3>
            <dl className="mt-4 space-y-4">
              <div>
                <dt className="text-sm font-medium text-gray-500">Package ID</dt>
                <dd className="mt-1 text-sm text-gray-900">{packageId}</dd>
              </div>
              {tracking_number && (
                <div>
                  <dt className="text-sm font-medium text-gray-500">Tracking Number</dt>
                  <dd className="mt-1 text-sm text-gray-900">{tracking_number}</dd>
                </div>
              )}
              {carrier && (
                <div>
                  <dt className="text-sm font-medium text-gray-500">Carrier</dt>
                  <dd className="mt-1 text-sm text-gray-900">{carrier}</dd>
                </div>
              )}
              {estimated_delivery && (
                <div>
                  <dt className="text-sm font-medium text-gray-500">Estimated Delivery</dt>
                  <dd className="mt-1 text-sm text-gray-900">
                    {format(new Date(estimated_delivery), 'MMM d, yyyy')}
                  </dd>
                </div>
              )}
            </dl>
          </div>
          <div>
            <h3 className="text-lg font-medium leading-6 text-gray-900">Current Status</h3>
            <div className="mt-4">
              <StatusBadge status={current_status.status} />
              <p className="mt-2 text-sm text-gray-500">{current_status.description}</p>
              <p className="mt-1 text-sm text-gray-900">{current_status.location}</p>
              <p className="mt-1 text-sm text-gray-500">
                {format(new Date(current_status.timestamp), 'MMM d, yyyy h:mm a')}
              </p>
            </div>
          </div>
        </div>
        <div className="mt-8">
          <h3 className="text-lg font-medium leading-6 text-gray-900">Tracking History</h3>
          <div className="mt-4">
            <StatusTimeline history={history} />
          </div>
        </div>
      </div>
    </div>
  );
}; 