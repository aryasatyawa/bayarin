import { useEffect, useRef } from 'react';
import { useLogout } from './useLogout';
import { storage } from '@/utils/storage';

const IDLE_TIMEOUT_MINUTES = 15;
const ABSOLUTE_TIMEOUT_HOURS = 12;

export const useIdleTimeout = () => {
    const { logout } = useLogout();
    const idleTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null);
    const absoluteTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null);
    const lastActivityRef = useRef<number>(Date.now());
    const loginTimeRef = useRef<number>(Date.now());

    const resetIdleTimer = () => {
        // Clear existing idle timer
        if (idleTimerRef.current) {
            clearTimeout(idleTimerRef.current);
        }

        lastActivityRef.current = Date.now();

        // Set new idle timer
        idleTimerRef.current = setTimeout(() => {
            // Check if user is still logged in
            const token = storage.getToken();
            if (token) {
                console.warn('Session expired due to inactivity');
                logout(true);
            }
        }, IDLE_TIMEOUT_MINUTES * 60 * 1000); // Convert minutes to milliseconds
    };

    const handleActivity = () => {
        const timeSinceLastActivity = (Date.now() - lastActivityRef.current) / 1000;

        // Only reset if at least 1 second has passed since last activity
        if (timeSinceLastActivity > 1) {
            resetIdleTimer();
        }
    };

    useEffect(() => {
        // Check if user is logged in
        const token = storage.getToken();
        if (!token) {
            return;
        }

        // Store login time
        loginTimeRef.current = Date.now();

        // Add event listeners for user activity
        const events = ['mousedown', 'keydown', 'scroll', 'touchstart', 'click'];

        events.forEach(event => {
            window.addEventListener(event, handleActivity);
        });

        // Initialize idle timer
        resetIdleTimer();

        // Set absolute timeout (12 hours)
        absoluteTimerRef.current = setTimeout(() => {
            console.warn('Session expired due to absolute timeout');
            logout(true);
        }, ABSOLUTE_TIMEOUT_HOURS * 60 * 60 * 1000); // Convert hours to milliseconds

        // Cleanup
        return () => {
            events.forEach(event => {
                window.removeEventListener(event, handleActivity);
            });

            if (idleTimerRef.current) {
                clearTimeout(idleTimerRef.current);
            }

            if (absoluteTimerRef.current) {
                clearTimeout(absoluteTimerRef.current);
            }
        };
    }, [logout]);
};
