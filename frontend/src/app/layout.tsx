import type { Metadata } from 'next';
import { Inter } from 'next/font/google';
import './globals.css';

const inter = Inter({ subsets: ['latin'] });

export const metadata: Metadata = {
  title: 'Moonlay Task Management',
  description: 'Aplikasi manajemen task sederhana dengan CRUD, autentikasi JWT, dan AI Chatbot berbasis RAG — Technical Test PT Moonlay Technologies',
  keywords: ['task management', 'moonlay', 'productivity'],
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="id" className="dark">
      <body className={`${inter.className} bg-slate-950 text-white antialiased`}>
        {children}
      </body>
    </html>
  );
}
