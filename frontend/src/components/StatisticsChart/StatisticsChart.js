import React from 'react';
import {
    Chart as ChartJS,
    CategoryScale,
    LinearScale,
    PointElement,
    LineElement,
    BarElement,
    Title,
    Tooltip,
    Legend,
    ArcElement
} from 'chart.js';
import { Line, Bar, Pie } from 'react-chartjs-2';

ChartJS.register(
    CategoryScale,
    LinearScale,
    PointElement,
    LineElement,
    BarElement,
    Title,
    Tooltip,
    Legend,
    ArcElement
);

export function StatisticsChart({ statistics }) {
    // Konfigurasi umum untuk semua chart
    const commonOptions = {
        responsive: true,
        plugins: {
            legend: {
                labels: {
                    color: '#333' // Warna teks legend selalu gelap
                }
            }
        },
        scales: {
            x: {
                grid: {
                    color: '#ddd' // Warna grid selalu terang
                },
                ticks: {
                    color: '#333' // Warna teks axis selalu gelap
                }
            },
            y: {
                grid: {
                    color: '#ddd' // Warna grid selalu terang
                },
                ticks: {
                    color: '#333' // Warna teks axis selalu gelap
                }
            }
        }
    };

    // Data untuk Line Chart
    const lineChartData = {
        labels: ['Min', 'Median', 'Mean', 'Max'],
        datasets: [{
            label: 'Konsumsi Energi',
            data: [
                statistics.minimum,
                statistics.median,
                statistics.rata_rata,
                statistics.maximum
            ],
            borderColor: 'rgb(75, 192, 192)',
            backgroundColor: 'rgba(75, 192, 192, 0.5)',
            tension: 0.1
        }]
    };

    // Data untuk Bar Chart
    const barChartData = {
        labels: ['Total Konsumsi', 'Rata-rata', 'Median'],
        datasets: [{
            label: 'Nilai Statistik',
            data: [
                statistics.total_konsumsi,
                statistics.rata_rata,
                statistics.median
            ],
            backgroundColor: [
                'rgba(255, 99, 132, 0.5)',
                'rgba(54, 162, 235, 0.5)',
                'rgba(255, 206, 86, 0.5)'
            ]
        }]
    };

    // Data untuk Pie Chart
    const pieChartData = {
        labels: ['Di bawah rata-rata', 'Di atas rata-rata'],
        datasets: [{
            data: [
                statistics.minimum,
                statistics.maximum - statistics.minimum
            ],
            backgroundColor: [
                'rgba(54, 162, 235, 0.5)',
                'rgba(255, 99, 132, 0.5)'
            ]
        }]
    };

    // Opsi khusus untuk Pie Chart
    const pieOptions = {
        plugins: {
            legend: {
                labels: {
                    color: '#333' // Warna teks legend selalu gelap
                }
            }
        }
    };

    return (
        <div className="charts-container">
            <div className="chart-item" style={{ background: 'white', padding: '20px', borderRadius: '10px' }}>
                <h3 style={{ color: '#333' }}>Tren Konsumsi Energi</h3>
                <Line data={lineChartData} options={commonOptions} />
            </div>
            <div className="chart-item" style={{ background: 'white', padding: '20px', borderRadius: '10px' }}>
                <h3 style={{ color: '#333' }}>Perbandingan Nilai Statistik</h3>
                <Bar data={barChartData} options={commonOptions} />
            </div>
            <div className="chart-item" style={{ background: 'white', padding: '20px', borderRadius: '10px' }}>
                <h3 style={{ color: '#333' }}>Distribusi Konsumsi</h3>
                <Pie data={pieChartData} options={pieOptions} />
            </div>
        </div>
    );
}
