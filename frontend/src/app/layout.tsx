import type { Metadata } from 'next';
import { Inter } from 'next/font/google';
import './globals.css';

const inter = Inter({ subsets: ['latin'] });

export const metadata: Metadata = {
  title: 'Moonlay Task Management',
  description: 'Aplikasi manajemen task sederhana dengan CRUD, autentikasi JWT, dan AI Chatbot berbasis RAG — Technical Test PT Moonlay Technologies',
  keywords: ['task management', 'moonlay', 'productivity'],
  viewport: {
    width: 'device-width',
    initialScale: 1,
    maximumScale: 1,
    userScalable: false,
  },
  colorScheme: 'dark',
  appleWebApp: {
    capable: true,
    statusBarStyle: 'black-translucent',
  },
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="id" className="dark scroll-smooth">
      <head>
        <meta name="theme-color" content="#0f0f0f" />
        <meta name="apple-mobile-web-app-capable" content="yes" />
        <meta name="apple-mobile-web-app-status-bar-style" content="black-translucent" />
      </head>
      <body className={`${inter.className} bg-slate-950 text-slate-50 antialiased overflow-x-hidden`}>
        {children}
      </body>
    </html>
  );
}