import { CreditCard } from 'lucide-react';
import styles from './DashboardHeader.module.css';

export const DashboardHeader = () => (
  <header className={styles.header}>
    <h1 className={styles.title}>
      <CreditCard className="text-blue-600 w-10 h-10" />
      PayFlux Dashboard
    </h1>
    <p className={styles.subtitle}>Manage payments, reconciliation, and risk evaluation in real-time.</p>
  </header>
);
