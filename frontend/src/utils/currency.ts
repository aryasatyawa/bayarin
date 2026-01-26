/**
 * Format amount from minor unit (integer) to IDR currency string
 * @param amount - Amount in minor unit (e.g., 10000000 = Rp 100.000)
 * @returns Formatted currency string
 */
export const formatCurrency = (amount: number): string => {
    // Convert from minor unit to major unit (divide by 100)
    const majorUnit = amount / 100;

    return new Intl.NumberFormat('id-ID', {
        style: 'currency',
        currency: 'IDR',
        minimumFractionDigits: 0,
        maximumFractionDigits: 0,
    }).format(majorUnit);
};

/**
 * Parse IDR currency string to minor unit (integer)
 * @param currencyString - Currency string (e.g., "Rp 100.000")
 * @returns Amount in minor unit
 */
export const parseCurrency = (currencyString: string): number => {
    // Remove all non-digit characters
    const cleanNumber = currencyString.replace(/[^\d]/g, '');

    // Convert to number and multiply by 100 (to minor unit)
    return parseInt(cleanNumber) * 100;
};

/**
 * Format amount from user input to minor unit
 * @param amount - Amount in major unit (e.g., 100000)
 * @returns Amount in minor unit
 */
export const toMinorUnit = (amount: number): number => {
    return Math.round(amount * 100);
};

/**
 * Format amount from minor unit to major unit
 * @param amount - Amount in minor unit
 * @returns Amount in major unit
 */
export const toMajorUnit = (amount: number): number => {
    return amount / 100;
};