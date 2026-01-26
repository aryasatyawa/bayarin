import { format, formatDistanceToNow } from 'date-fns';
import { id } from 'date-fns/locale';

/**
 * Format date to readable format
 */
export const formatDate = (date: string | Date, formatStr: string = 'dd MMM yyyy'): string => {
    return format(new Date(date), formatStr, { locale: id });
};

/**
 * Format datetime to readable format
 */
export const formatDateTime = (date: string | Date): string => {
    return format(new Date(date), 'dd MMM yyyy HH:mm', { locale: id });
};

/**
 * Format relative time (e.g., "2 hours ago")
 */
export const formatRelativeTime = (date: string | Date): string => {
    return formatDistanceToNow(new Date(date), { addSuffix: true, locale: id });
};