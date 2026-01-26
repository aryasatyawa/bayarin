import React from 'react';

interface CardProps {
    children: React.ReactNode;
    className?: string;
    onClick?: () => void;
    variant?: 'default' | 'elevated' | 'interactive';
}

export const Card: React.FC<CardProps> = ({
    children,
    className = '',
    onClick,
    variant = 'default'
}) => {
    const variants = {
        default: 'bg-white border border-gray-200 shadow-sm hover:shadow-md',
        elevated: 'bg-white border border-gray-200 shadow-md hover:shadow-lg',
        interactive: 'bg-white border border-gray-200 shadow-sm hover:shadow-xl hover:scale-[1.02] cursor-pointer active:scale-[0.98]',
    };

    return (
        <div
            className={`card transition-all duration-300 ${variants[variant]} ${className}`}
            onClick={onClick}
        >
            {children}
        </div>
    );
};