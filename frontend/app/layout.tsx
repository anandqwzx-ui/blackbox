import type { Metadata } from 'next';
import './globals.css';

export const metadata: Metadata = {
  title: 'Crudbox',
  description: 'Create and manage mock APIs for rapid prototyping.',
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <body className="bg-black text-white min-h-screen antialiased">
        <main className="max-w-5xl mx-auto px-4 py-10">{children}</main>
      </body>
    </html>
  );
}
