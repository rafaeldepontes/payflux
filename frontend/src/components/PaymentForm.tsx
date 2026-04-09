import { useState } from 'react';
import axios from 'axios';
import { ArrowRightLeft, History } from 'lucide-react';
import { generateIdempotencyKey } from '../utils/helpers';
import type { PaymentRes } from '../types';
import styles from './PaymentForm.module.css';

const LEDGER_URL = import.meta.env.VITE_LEDGER_URL || 'http://localhost:8080';

interface Props {
  onPaymentCreated: (payment: PaymentRes) => void;
}

export const PaymentForm = ({ onPaymentCreated }: Props) => {
  const [fromAccount, setFromAccount] = useState('1');
  const [toAccount, setToAccount] = useState('2');
  const [amount, setAmount] = useState('100');
  const [idempotencyKey, setIdempotencyKey] = useState(generateIdempotencyKey());
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async () => {
    setLoading(true);
    setError('');
    try {
      const res = await axios.post(`${LEDGER_URL}/payments`, {
        from_account: parseInt(fromAccount),
        to_account: parseInt(toAccount),
        amount: parseInt(amount),
        currency: 'USD'
      }, {
        headers: { 'Idempotency-Key': idempotencyKey }
      });
      onPaymentCreated(res.data);
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to create payment');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className={styles.card}>
      <h2 className={styles.title}>
        <ArrowRightLeft className="w-5 h-5 text-blue-500" />
        Create Payment
      </h2>
      <div className="space-y-4">
        <div className={styles.formGroup}>
          <label className={styles.label}>From Account</label>
          <input type="number" value={fromAccount} onChange={e => setFromAccount(e.target.value)} className={styles.input} />
        </div>
        <div className={styles.formGroup}>
          <label className={styles.label}>To Account</label>
          <input type="number" value={toAccount} onChange={e => setToAccount(e.target.value)} className={styles.input} />
        </div>
        <div className={styles.formGroup}>
          <label className={styles.label}>Amount</label>
          <input type="number" value={amount} onChange={e => setAmount(e.target.value)} className={styles.input} />
        </div>
        <div className={styles.formGroup}>
          <label className={styles.label}>Idempotency Key</label>
          <div className="flex gap-2">
            <input type="text" value={idempotencyKey} onChange={e => setIdempotencyKey(e.target.value)} className={`${styles.input} text-xs font-mono`} />
            <button onClick={() => setIdempotencyKey(generateIdempotencyKey())} className="p-2 text-slate-400 hover:text-blue-500">
              <History className="w-4 h-4" />
            </button>
          </div>
        </div>
        <button onClick={handleSubmit} disabled={loading} className={styles.button}>
          {loading ? 'Processing...' : 'Send Payment'}
        </button>
        {error && <p className="text-red-500 text-sm mt-2">{error}</p>}
      </div>
    </div>
  );
};
