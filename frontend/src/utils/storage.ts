export const storage = {
    // Save token
    setToken: (token: string): void => {
        localStorage.setItem('token', token);
    },

    // Get token
    getToken: (): string | null => {
        return localStorage.getItem('token');
    },

    // Remove token
    removeToken: (): void => {
        localStorage.removeItem('token');
    },

    // Save user data
    setUser: (user: any): void => {
        localStorage.setItem('user', JSON.stringify(user));
    },

    // Get user data
    getUser: (): any | null => {
        const user = localStorage.getItem('user');
        return user ? JSON.parse(user) : null;
    },

    // Remove user data
    removeUser: (): void => {
        localStorage.removeItem('user');
    },

    // Clear all storage
    clear: (): void => {
        localStorage.clear();
    },
};