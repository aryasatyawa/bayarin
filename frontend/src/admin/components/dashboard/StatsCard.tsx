import React from 'react';
import { LucideIcon } from 'lucide-react';

interface StatsCardProps {
    title: string;
    value: string | number;
    icon: LucideIcon;
    trend?: {
        value: string;
        isPositive: boolean;
    };
    color: 'blue' | 'green' | 'yellow' | 'red' | 'purple';
}

const colorClasses = {
    blue: 'bg-blue-500/10 text-blue-500 border-blue-500/20',
    green: 'bg-green-500/10 text-green-500 border-green-500/20',
    yellow: 'bg-yellow-500/10 text-yellow-500 border-yellow-500/20',
    red: 'bg-red-500/10 text-red-500 border-red-500/20',
    purple: 'bg-purple-500/10 text-purple-500 border-purple-500/20',
};

export const StatsCard: React.FC<StatsCardProps> = ({
    title,
    value,
    icon: Icon,
    trend,
    color,
}) => {
    return (
        <div className="bg-white rounded-xl p-6 border border-gray-200 shadow-sm hover:border-gray-300 transition-all">
            <div className="flex items-start justify-between">
                <div className="flex-1">
                    <p className="text-gray-500 text-sm font-medium mb-2">{title}</p>
                    <h3 className="text-gray-900 text-3xl font-bold mb-2">{value}</h3>
                    {trend && (
                        <p className={`text-sm ${trend.isPositive ? 'text-green-500' : 'text-red-500'}`}>
                            {trend.value}
                        </p>
                    )}
                </div>
                <div className={`flex items-center justify-center shrink-0 w-12 h-12 rounded-lg border ${colorClasses[color]}`}>
                    <Icon className="w-6 h-6" strokeWidth={2} />
                </div>
            </div>
        </div>
    );
};