import React from 'react';
import ReactDOM from 'react-dom/client';
import { ThemeProvider } from '@/components/ThemeProvider.tsx';
import { Toaster } from '@/components/ui/toaster';
import App from './App.tsx';
import './index.css';

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <ThemeProvider defaultTheme="dark" storageKey="task-manager-theme">
      <App />
      <Toaster />
    </ThemeProvider>
  </React.StrictMode>
);
