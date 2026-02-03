import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { motion, AnimatePresence, Variants } from 'framer-motion';
import { Button } from '../components/common/Button';
import {
    ShieldCheck,
    Zap,
    QrCode,
    CreditCard,
    Menu,
    X,
    Smartphone,
    Globe,
    Lock,
    Server,
    CheckCircle2,
    ArrowRight,
    Bell
} from 'lucide-react';

// --- TYPES & DATA ---

interface FeatureItem {
    icon: React.ElementType;
    title: string;
    desc: string;
}

interface StepItem {
    number: string;
    title: string;
    desc: string;
}

const FEATURES: FeatureItem[] = [
    {
        icon: Zap,
        title: 'Fast Transactions',
        desc: 'Settlement real-time di bawah 2 detik dengan infrastruktur high-availability.',
    },
    {
        icon: Server,
        title: 'Secure Ledger',
        desc: 'Sistem pencatatan immutable yang menjamin integritas data transaksi.',
    },
    {
        icon: QrCode,
        title: 'QR Transfer',
        desc: 'Mendukung standar QRIS dan pembayaran lintas border secara instan.',
    },
    {
        icon: Globe,
        title: 'Gateway Ready',
        desc: 'API yang mudah diintegrasikan untuk e-commerce dan platform SaaS.',
    },
];

const STEPS: StepItem[] = [
    {
        number: '01',
        title: 'Create Wallet',
        desc: 'Daftar dalam hitungan menit dengan verifikasi e-KYC otomatis.',
    },
    {
        number: '02',
        title: 'Scan or Input',
        desc: 'Scan QRIS merchant atau input nomor tujuan transfer dengan mudah.',
    },
    {
        number: '03',
        title: 'Payment Completed',
        desc: 'Transaksi selesai, notifikasi instan, dan resi digital tersimpan aman.',
    },
];

// --- ANIMATION VARIANTS ---

const EASE_OUT: [number, number, number, number] = [0.16, 1, 0.3, 1];

const fadeInUp: Variants = {
    hidden: { opacity: 0, y: 40 },
    visible: {
        opacity: 1,
        y: 0,
        transition: { duration: 0.6, ease: EASE_OUT },
    },
};

const staggerContainer: Variants = {
    hidden: { opacity: 0 },
    visible: {
        opacity: 1,
        transition: {
            staggerChildren: 0.2,
        },
    },
};

// --- COMPONENTS ---

const Navbar = () => {
    const [isScrolled, setIsScrolled] = useState(false);
    const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);

    useEffect(() => {
        const handleScroll = () => setIsScrolled(window.scrollY > 20);
        window.addEventListener('scroll', handleScroll);
        return () => window.removeEventListener('scroll', handleScroll);
    }, []);

    const navLinks = ['Features', 'Security', 'How It Works', 'Contact'];

    return (
        <nav
            className={`fixed top-0 left-0 right-0 z-50 transition-all duration-300 ${isScrolled ? 'bg-white/80 backdrop-blur-md shadow-sm py-4' : 'bg-transparent py-6'
                }`}
        >
            <div className="container mx-auto px-6 flex justify-between items-center">
                <div className="font-bold text-2xl text-slate-900 tracking-tighter cursor-pointer">
                    Bayarin<span className="text-blue-600">.</span>
                </div>

                {/* Desktop Menu */}
                <div className="hidden md:flex gap-8 items-center">
                    {navLinks.map((link) => (
                        <a
                            key={link}
                            href={`#${link.toLowerCase().replace(/\s/g, '-')}`}
                            className="text-slate-600 hover:text-blue-600 text-sm font-medium transition-colors"
                        >
                            {link}
                        </a>
                    ))}
                    <Link
                        to="/login"
                        className="px-5 py-2 bg-slate-900 text-white text-sm font-medium rounded-full hover:bg-slate-800 transition-transform hover:scale-105"
                    >
                        Login
                    </Link>
                </div>

                {/* Mobile Toggle */}
                <button
                    className="md:hidden text-slate-900"
                    onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
                >
                    {isMobileMenuOpen ? <X size={24} /> : <Menu size={24} />}
                </button>
            </div>

            {/* Mobile Menu */}
            <AnimatePresence>
                {isMobileMenuOpen && (
                    <motion.div
                        initial={{ opacity: 0, height: 0 }}
                        animate={{ opacity: 1, height: 'auto' }}
                        exit={{ opacity: 0, height: 0 }}
                        className="md:hidden bg-white border-b border-slate-100 overflow-hidden"
                    >
                        <div className="flex flex-col p-6 gap-4">
                            {navLinks.map((link) => (
                                <a
                                    key={link}
                                    href={`#${link.toLowerCase().replace(/\s/g, '-')}`}
                                    className="text-slate-600 font-medium"
                                    onClick={() => setIsMobileMenuOpen(false)}
                                >
                                    {link}
                                </a>
                            ))}
                            <Link
                                to="/login"
                                className="w-full py-3 bg-slate-100 text-slate-800 rounded-lg font-medium text-center"
                                onClick={() => setIsMobileMenuOpen(false)}
                            >
                                Login
                            </Link>
                            <Link
                                to="/register"
                                className="w-full py-3 bg-blue-600 text-white rounded-lg font-medium text-center"
                                onClick={() => setIsMobileMenuOpen(false)}
                            >
                                Get Started
                            </Link>
                        </div>
                    </motion.div>
                )}
            </AnimatePresence>
        </nav>
    );
};

export function Hero() {
    const containerVariants: Variants = {
        hidden: { opacity: 0 },
        visible: {
            opacity: 1,
            transition: {
                staggerChildren: 0.15,
                delayChildren: 0.2,
            },
        },
    };

    const itemVariants: Variants = {
        hidden: { opacity: 0, y: 20 },
        visible: {
            opacity: 1,
            y: 0,
            transition: {
                duration: 0.6,
                ease: EASE_OUT,
            },
        },
    };

    const floatingAnimation = {
        y: [-10, 10, -10],
        transition: {
            duration: 4,
            repeat: Infinity,
            ease: EASE_OUT,
        },
    };

    return (
        <section className="relative pt-24 pb-16 lg:pt-24 lg:pb-10 overflow-hidden">
            {/* Background Blobs */}
            <div className="absolute top-0 right-0 -z-10 w-[600px] h-[600px] bg-blue-50 rounded-full blur-3xl opacity-50 translate-x-1/3 -translate-y-1/4" />
            <div className="absolute bottom-0 left-0 -z-10 w-[400px] h-[400px] bg-indigo-50 rounded-full blur-3xl opacity-50 -translate-x-1/3 translate-y-1/4" />

            <div className="container mx-auto px-4 md:px-6">
                <div className="grid lg:grid-cols-2 gap-12 lg:gap-8 items-center">
                    {/* Left Content */}
                    <motion.div
                        variants={containerVariants}
                        initial="hidden"
                        animate="visible"
                        className="max-w-2xl"
                    >
                        <motion.div
                            variants={itemVariants}
                            className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-blue-50 text-blue-600 text-sm font-medium mb-6"
                        >
                            <span className="relative flex h-2 w-2">
                                <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-blue-400 opacity-75" />
                                <span className="relative inline-flex rounded-full h-2 w-2 bg-blue-500" />
                            </span>
                            Now available in 30+ countries
                        </motion.div>

                        <motion.h1
                            variants={itemVariants}
                            className="text-5xl md:text-6xl lg:text-7xl font-bold tracking-tight text-slate-900 mb-6 leading-[1.1]"
                        >
                            Your Digital Wallet,{' '}
                            <span className="text-blue-600">Simplified.</span>
                        </motion.h1>

                        <motion.p
                            variants={itemVariants}
                            className="text-lg md:text-xl text-slate-600 mb-8 max-w-lg leading-relaxed"
                        >
                            Send money, pay bills, and track your spending with the most secure
                            and intuitive e-wallet experience designed for modern life.
                        </motion.p>

                        <motion.div
                            variants={itemVariants}
                            className="flex flex-col sm:flex-row gap-4"
                        >
                            <Button size="lg" className="group">
                                Get Started
                                <ArrowRight className="ml-2 h-4 w-4 transition-transform group-hover:translate-x-1" />
                            </Button>

                            <Button size="lg" variant="secondary">
                                Learn More
                            </Button>
                        </motion.div>

                        <motion.div
                            variants={itemVariants}
                            className="mt-10 flex items-center gap-6 text-sm text-slate-500"
                        >
                            <div className="flex items-center gap-2">
                                <ShieldCheck className="h-5 w-5 text-blue-600" />
                                <span>Bank-grade security</span>
                            </div>
                            <div className="flex items-center gap-2">
                                <CreditCard className="h-5 w-5 text-blue-600" />
                                <span>No hidden fees</span>
                            </div>
                        </motion.div>
                    </motion.div>

                    {/* Right Content – Phone Mockup */}
                    <motion.div
                        initial={{ opacity: 0, x: 50 }}
                        animate={{ opacity: 1, x: 0 }}
                        transition={{ duration: 0.8, delay: 0.4 }}
                        className="relative lg:h-[600px] flex items-center justify-center"
                    >
                        {/* Phone Body */}
                        <div className="relative w-[300px] h-[600px] bg-slate-900 rounded-[3rem] border-8 border-slate-900 shadow-2xl overflow-hidden z-10">
                            <div className="absolute inset-0 bg-white flex flex-col">
                                {/* Status Bar */}
                                <div className="h-6 bg-slate-900 flex justify-center items-center">
                                    <div className="w-12 h-3 bg-black rounded-full" />
                                </div>

                                {/* Header */}
                                <div className="p-6 bg-blue-600 text-white pb-8 rounded-b-[2rem]">
                                    <div className="flex justify-between mb-6">
                                        <div className="w-8 h-8 rounded-full bg-white/20" />
                                        <div className="w-8 h-8 rounded-full bg-white/20" />
                                    </div>
                                    <div className="text-sm opacity-80">Total Balance</div>
                                    <div className="text-3xl font-bold mt-1">$12,450.00</div>
                                </div>

                                {/* Body */}
                                <div className="p-6 flex-1 bg-slate-50">
                                    <div className="flex justify-between gap-4 mb-6">
                                        {['Send', 'Request', 'Top Up', 'More'].map((label) => (
                                            <div key={label} className="flex flex-col items-center gap-2">
                                                <div className="w-12 h-12 bg-white rounded-full shadow-sm flex items-center justify-center">
                                                    <div className="w-5 h-5 bg-blue-100 rounded-full" />
                                                </div>
                                                <span className="text-[10px] font-medium text-slate-600">
                                                    {label}
                                                </span>
                                            </div>
                                        ))}
                                    </div>

                                    <div className="space-y-3">
                                        <div className="text-xs font-semibold text-slate-400 uppercase tracking-wider">
                                            Recent Activity
                                        </div>

                                        {[
                                            { name: 'Netflix', amount: '14.99', color: 'orange' },
                                            { name: 'Spotify', amount: '9.99', color: 'green' },
                                            { name: 'Amazon', amount: '45.50', color: 'purple' },
                                        ].map((item) => (
                                            <div
                                                key={item.name}
                                                className="bg-white p-3 rounded-xl shadow-sm flex justify-between items-center"
                                            >
                                                <div className="flex items-center gap-3">
                                                    <div
                                                        className={`w-10 h-10 rounded-full bg-${item.color}-100 text-${item.color}-600 flex items-center justify-center font-bold`}
                                                    >
                                                        {item.name[0]}
                                                    </div>
                                                    <div>
                                                        <div className="text-sm font-bold text-slate-800">
                                                            {item.name}
                                                        </div>
                                                        <div className="text-xs text-slate-500">
                                                            Subscription
                                                        </div>
                                                    </div>
                                                </div>
                                                <div className="text-sm font-bold text-slate-900">
                                                    -${item.amount}
                                                </div>
                                            </div>
                                        ))}
                                    </div>
                                </div>
                            </div>
                        </div>

                        {/* Floating Cards */}
                        <motion.div
                            animate={floatingAnimation}
                            className="absolute top-20 -left-12 bg-white p-4 rounded-2xl shadow-xl z-20 flex items-center gap-3"
                        >
                            <div className="w-10 h-10 bg-green-100 rounded-full flex items-center justify-center text-green-600">
                                <ArrowRight className="w-5 h-5 -rotate-45" />
                            </div>
                            <div>
                                <div className="text-xs text-slate-500">Payment Received</div>
                                <div className="text-sm font-bold text-slate-900">+$250.00</div>
                            </div>
                        </motion.div>

                        <motion.div
                            animate={{
                                y: [10, -10, 10],
                                transition: {
                                    duration: 5,
                                    repeat: Infinity,
                                    ease: EASE_OUT,
                                    delay: 1,
                                },
                            }}
                            className="absolute bottom-32 -right-8 bg-white p-4 rounded-2xl shadow-xl z-20 flex items-center gap-3"
                        >
                            <div className="w-10 h-10 bg-blue-100 rounded-full flex items-center justify-center text-blue-600">
                                <Bell className="w-5 h-5" />
                            </div>
                            <div>
                                <div className="text-xs text-slate-500">New Reward</div>
                                <div className="text-sm font-bold text-slate-900">
                                    5% Cashback
                                </div>
                            </div>
                        </motion.div>
                    </motion.div>
                </div>
            </div>
        </section>
    );
}

const Features = () => {
    return (
        <section id="features" className="py-20 bg-white relative">
            <div className="container mx-auto px-6">
                <div className="text-center mb-16">
                    <h2 className="text-3xl md:text-4xl font-bold text-slate-900 mb-4">Why Choose Bayarin?</h2>
                    <p className="text-slate-600 max-w-2xl mx-auto">
                        Infrastruktur finansial yang dibangun untuk skalabilitas dan kecepatan tanpa mengorbankan
                        keamanan.
                    </p>
                </div>

                <motion.div
                    variants={staggerContainer}
                    initial="hidden"
                    whileInView="visible"
                    viewport={{ once: true, margin: '-50px' }}
                    className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8"
                >
                    {FEATURES.map((feature, idx) => (
                        <motion.div
                            key={idx}
                            variants={fadeInUp}
                            whileHover={{ y: -5 }}
                            className="p-8 rounded-2xl bg-slate-50 hover:bg-white hover:shadow-xl border border-slate-100 hover:border-blue-100 transition-all duration-300 group"
                        >
                            <div className="w-12 h-12 rounded-lg bg-blue-100 text-blue-600 flex items-center justify-center mb-6 group-hover:bg-blue-600 group-hover:text-white transition-colors">
                                <feature.icon size={24} />
                            </div>
                            <h3 className="text-xl font-bold text-slate-900 mb-3">{feature.title}</h3>
                            <p className="text-slate-600 text-sm leading-relaxed">{feature.desc}</p>
                        </motion.div>
                    ))}
                </motion.div>
            </div>
        </section>
    );
};

const HowItWorks = () => {
    return (
        <section id="how-it-works" className="py-24 bg-slate-50 overflow-hidden">
            <div className="container mx-auto px-6">
                <div className="flex flex-col lg:flex-row items-center gap-16">
                    <div className="lg:w-1/2">
                        <h2 className="text-3xl md:text-4xl font-bold text-slate-900 mb-6">
                            Transaksi Semudah <br /> Menjentikkan Jari
                        </h2>
                        <p className="text-slate-600 mb-10 text-lg">
                            Kami menyederhanakan proses finansial yang kompleks menjadi pengalaman pengguna yang
                            intuitif. Tidak ada lagi form manual yang panjang.
                        </p>

                        <div className="space-y-8">
                            {STEPS.map((step, idx) => (
                                <motion.div
                                    initial={{ opacity: 0, x: -20 }}
                                    whileInView={{ opacity: 1, x: 0 }}
                                    viewport={{ once: true }}
                                    transition={{ delay: idx * 0.2 }}
                                    key={idx}
                                    className="flex gap-6"
                                >
                                    <div className="flex-shrink-0 w-12 h-12 rounded-full bg-white border border-blue-200 flex items-center justify-center text-blue-600 font-bold shadow-sm">
                                        {step.number}
                                    </div>
                                    <div>
                                        <h4 className="text-xl font-bold text-slate-900 mb-2">{step.title}</h4>
                                        <p className="text-slate-600 text-sm">{step.desc}</p>
                                    </div>
                                </motion.div>
                            ))}
                        </div>
                    </div>

                    <div className="lg:w-1/2 relative">
                        <motion.div
                            initial={{ opacity: 0, scale: 0.9 }}
                            whileInView={{ opacity: 1, scale: 1 }}
                            viewport={{ once: true }}
                            transition={{ duration: 0.5 }}
                            className="bg-white p-8 rounded-3xl shadow-2xl border border-slate-100 relative z-10"
                        >
                            <div className="flex justify-between items-center mb-8 border-b border-slate-100 pb-4">
                                <div>
                                    <p className="text-xs text-slate-500 uppercase tracking-wide mb-1">Total Balance</p>
                                    <h3 className="text-3xl font-bold text-slate-900">Rp 24.500.000</h3>
                                </div>
                                <div className="w-10 h-10 bg-blue-50 rounded-full flex items-center justify-center text-blue-600">
                                    <Smartphone size={20} />
                                </div>
                            </div>
                            <div className="space-y-4">
                                {[1, 2, 3].map((_, i) => (
                                    <div
                                        key={i}
                                        className="flex items-center justify-between p-3 rounded-xl bg-slate-50 hover:bg-slate-100 transition-colors"
                                    >
                                        <div className="flex items-center gap-3">
                                            <div className="w-10 h-10 rounded-full bg-slate-200"></div>
                                            <div>
                                                <div className="w-24 h-3 bg-slate-300 rounded mb-2"></div>
                                                <div className="w-16 h-2 bg-slate-200 rounded"></div>
                                            </div>
                                        </div>
                                        <div className="w-12 h-4 bg-slate-200 rounded"></div>
                                    </div>
                                ))}
                            </div>

                            {/* Floating Elements for Decor */}
                            <motion.div
                                animate={{ y: [0, -10, 0] }}
                                transition={{ repeat: Infinity, duration: 4, ease: 'easeInOut' }}
                                className="absolute -top-6 -right-6 bg-blue-600 text-white p-4 rounded-xl shadow-lg"
                            >
                                <CheckCircle2 size={24} />
                            </motion.div>
                        </motion.div>
                        <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[120%] h-[120%] bg-blue-100/50 rounded-full blur-3xl -z-10"></div>
                    </div>
                </div>
            </div>
        </section>
    );
};

const Security = () => {
    return (
        <section id="security" className="py-24 bg-white">
            <div className="container mx-auto px-6">
                <div className="bg-slate-900 rounded-3xl p-8 md:p-16 text-white overflow-hidden relative">
                    <div className="absolute top-0 right-0 w-96 h-96 bg-blue-600/20 rounded-full blur-3xl -translate-y-1/2 translate-x-1/2"></div>

                    <div className="grid md:grid-cols-2 gap-12 items-center relative z-10">
                        <div>
                            <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-blue-900/50 border border-blue-700 text-blue-300 text-xs font-semibold uppercase tracking-wider mb-6">
                                <ShieldCheck size={14} />
                                Bank-Grade Security
                            </div>
                            <h2 className="text-3xl md:text-4xl font-bold mb-6">Keamanan Adalah Prioritas Mutlak</h2>
                            <p className="text-slate-400 mb-8 leading-relaxed">
                                Kami menggunakan teknologi enkripsi terkini dan kepatuhan regulasi ketat untuk memastikan
                                setiap bit data dan dana Anda terlindungi.
                            </p>

                            <div className="grid grid-cols-1 sm:grid-cols-2 gap-6">
                                <div className="flex gap-4">
                                    <div className="w-10 h-10 rounded-lg bg-blue-500/20 flex items-center justify-center text-blue-400 flex-shrink-0">
                                        <Lock size={20} />
                                    </div>
                                    <div>
                                        <h4 className="font-bold text-lg mb-1">ACID Transaction</h4>
                                        <p className="text-sm text-slate-400">Menjamin integritas data 100%.</p>
                                    </div>
                                </div>
                                <div className="flex gap-4">
                                    <div className="w-10 h-10 rounded-lg bg-blue-500/20 flex items-center justify-center text-blue-400 flex-shrink-0">
                                        <CreditCard size={20} />
                                    </div>
                                    <div>
                                        <h4 className="font-bold text-lg mb-1">PCI DSS Ready</h4>
                                        <p className="text-sm text-slate-400">Standar keamanan kartu pembayaran.</p>
                                    </div>
                                </div>
                            </div>
                        </div>

                        <div className="relative">
                            <div className="bg-slate-800 rounded-xl p-6 border border-slate-700 font-mono text-sm text-blue-300 shadow-2xl">
                                <div className="flex gap-2 mb-4 border-b border-slate-700 pb-2">
                                    <div className="w-3 h-3 rounded-full bg-red-500"></div>
                                    <div className="w-3 h-3 rounded-full bg-yellow-500"></div>
                                    <div className="w-3 h-3 rounded-full bg-green-500"></div>
                                </div>
                                <p className="mb-2">
                                    <span className="text-purple-400">const</span> transaction ={' '}
                                    <span className="text-yellow-300">await</span> ledger.
                                    <span className="text-blue-400">verify</span>(payload);
                                </p>
                                <p className="mb-2">
                                    <span className="text-purple-400">if</span> (!transaction.isValid){' '}
                                    <span className="text-purple-400">throw new</span> Error(
                                    <span className="text-green-300">'Fraud detected'</span>);
                                </p>
                                <p className="mb-2 text-slate-500">// Transaction encrypted with AES-256</p>
                                <p>
                                    <span className="text-purple-400">return</span>{' '}
                                    <span className="text-blue-400">finalize</span>(transaction);
                                </p>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </section>
    );
};

const CTA = () => {
    return (
        <section className="py-24 bg-blue-600 relative overflow-hidden">
            <div className="absolute inset-0 bg-[url('https://www.transparenttextures.com/patterns/cubes.png')] opacity-10"></div>
            <div className="container mx-auto px-6 relative z-10 text-center">
                <h2 className="text-4xl md:text-5xl font-bold text-white mb-6">
                    Siap Mengubah Cara Anda Bertransaksi?
                </h2>
                <p className="text-blue-100 text-lg mb-10 max-w-2xl mx-auto">
                    Bergabung dengan ribuan bisnis dan individu yang telah mempercayakan pengelolaan finansial
                    mereka pada Bayarin.
                </p>
                <motion.div whileHover={{ scale: 1.05 }} whileTap={{ scale: 0.95 }}>
                    <Link
                        to="/register"
                        className="inline-block bg-white text-blue-700 px-10 py-4 rounded-xl font-bold text-lg shadow-xl hover:shadow-2xl transition-all"
                    >
                        Try Bayarin Now
                    </Link>
                </motion.div>
            </div>
        </section>
    );
};

const Footer = () => {
    return (
        <footer id="contact" className="bg-slate-50 pt-20 pb-10 border-t border-slate-200">
            <div className="container mx-auto px-6">
                <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-5 gap-8 mb-16">
                    <div className="col-span-2 lg:col-span-2">
                        <div className="font-bold text-2xl text-slate-900 tracking-tighter mb-4">
                            Bayarin<span className="text-blue-600">.</span>
                        </div>
                        <p className="text-slate-500 text-sm leading-relaxed max-w-xs mb-6">
                            Platform pembayaran digital terdepan yang mengutamakan kecepatan, keamanan, dan kemudahan
                            integrasi.
                        </p>
                        <div className="flex gap-4 text-slate-400">
                            {/* Social Icons Placeholder */}
                            <a href="#" className="w-8 h-8 bg-white border border-slate-200 rounded-full flex items-center justify-center hover:border-blue-600 hover:text-blue-600 transition-colors cursor-pointer">
                                <Globe size={14} />
                            </a>
                            <a href="#" className="w-8 h-8 bg-white border border-slate-200 rounded-full flex items-center justify-center hover:border-blue-600 hover:text-blue-600 transition-colors cursor-pointer">
                                <Zap size={14} />
                            </a>
                        </div>
                    </div>

                    <div>
                        <h4 className="font-bold text-slate-900 mb-4">Product</h4>
                        <ul className="space-y-2 text-sm text-slate-600">
                            <li>
                                <a href="#features" className="hover:text-blue-600">
                                    Features
                                </a>
                            </li>
                            <li>
                                <a href="#" className="hover:text-blue-600">
                                    Pricing
                                </a>
                            </li>
                            <li>
                                <a href="#" className="hover:text-blue-600">
                                    API Docs
                                </a>
                            </li>
                        </ul>
                    </div>

                    <div>
                        <h4 className="font-bold text-slate-900 mb-4">Company</h4>
                        <ul className="space-y-2 text-sm text-slate-600">
                            <li>
                                <a href="#" className="hover:text-blue-600">
                                    About Us
                                </a>
                            </li>
                            <li>
                                <a href="#" className="hover:text-blue-600">
                                    Careers
                                </a>
                            </li>
                            <li>
                                <a href="#" className="hover:text-blue-600">
                                    Contact
                                </a>
                            </li>
                        </ul>
                    </div>

                    <div>
                        <h4 className="font-bold text-slate-900 mb-4">Legal</h4>
                        <ul className="space-y-2 text-sm text-slate-600">
                            <li>
                                <a href="#" className="hover:text-blue-600">
                                    Privacy Policy
                                </a>
                            </li>
                            <li>
                                <a href="#" className="hover:text-blue-600">
                                    Terms of Service
                                </a>
                            </li>
                        </ul>
                    </div>
                </div>

                <div className="border-t border-slate-200 pt-8 flex flex-col md:flex-row justify-between items-center gap-4">
                    <p className="text-slate-400 text-xs">© 2026 Bayarin Technologies. All rights reserved.</p>
                    <div className="flex gap-6 text-xs text-slate-500">
                        <span>Jakarta, Indonesia</span>
                    </div>
                </div>
            </div>
        </footer>
    );
};

export const InformationPage: React.FC = () => {
    return (
        <main className="font-sans antialiased text-slate-900 bg-white selection:bg-blue-100 selection:text-blue-900">
            <Navbar />
            <Hero />
            <Features />
            <HowItWorks />
            <Security />
            <CTA />
            <Footer />
        </main>
    );
};
