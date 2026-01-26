import { useEffect } from 'react';
import { useQueryClient } from '@tanstack/react-query';

export const useLogoutListener = () => {
    const queryClient = useQueryClient();

    useEffect(() => {
        const handleUnauthorizedLogout = () => {
            // Clear React Query cache
            queryClient.clear();
        };

        window.addEventListener('unauthorized-logout', handleUnauthorizedLogout);

        return () => {
            window.removeEventListener('unauthorized-logout', handleUnauthorizedLogout);
        };
    }, [queryClient]);
};
