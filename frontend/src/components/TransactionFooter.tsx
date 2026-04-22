import type { PaymentRes } from '../types';
import styles from './TransactionFooter.module.css';

interface Props {
  payment: PaymentRes | null;
}

export const TransactionFooter = ({ payment }: Props) => {
  if (!payment) return null;

  return (
    <div className={styles.footer}>
      <div>
        <p className={styles.label}>Last Transaction Successful</p>
        <p className={styles.txId}>ID: {payment.payment_id}</p>
      </div>
      <div>
        <p className={styles.amount}>${(payment.amount / 100).toFixed(2)}</p>
        <p className={styles.status}>{payment.status}</p>
      </div>
    </div>
  );
};
