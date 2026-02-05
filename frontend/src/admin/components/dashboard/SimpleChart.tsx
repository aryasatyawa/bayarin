import React from 'react';
import { formatCurrency } from '@/utils/currency';

interface ChartData {
    date: string;
    count: number;
    volume: number;
}

interface SimpleChartProps {
    data: ChartData[];
    title: string;
}

export const SimpleChart: React.FC<SimpleChartProps> = ({ data, title }) => {
    const maxVolume = Math.max(...data.map((d) => d.volume));

    return (
        <div className="bg-white rounded-xl p-6 border border-gray-200 shadow-sm">
            <h3 className="text-gray-900 font-semibold text-lg mb-6">{title}</h3>

            <div className="space-y-4">
                {data.slice(0, 7).map((item, index) => {
                    const percentage = (item.volume / maxVolume) * 100;

                    return (
                        <div key={index}>
                            <div className="flex justify-between text-sm mb-2">
                                <span className="text-gray-500">
                                    {new Date(item.date).toLocaleDateString('id-ID', {
                                        day: 'numeric',
                                        month: 'short',
                                    })}
                                </span>
                                <span className="text-gray-900 font-medium">
                                    {formatCurrency(item.volume)}
                                </span>
                            </div>
                            <div className="w-full bg-gray-200 rounded-full h-2">
                                <div
                                    className="bg-blue-500 h-2 rounded-full transition-all"
                                    style={{ width: `${percentage}%` }}
                                />
                            </div>
                            <p className="text-xs text-gray-500 mt-1">
                                {item.count} transaksi
                            </p>
                        </div>
                    );
                })}
            </div>
        </div>
    );
};  