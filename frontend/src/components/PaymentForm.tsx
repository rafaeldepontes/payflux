import { useRef, useState } from 'react';
import axios, { AxiosError } from 'axios';
import { ArrowRightLeft, History } from 'lucide-react';
import { formatWithCursor, generateIdempotencyKey, parseMoneyToCents } from '../utils/helpers';
import type { PaymentRes } from '../types';
import styles from './PaymentForm.module.css';

const LEDGER_URL = import.meta.env.VITE_LEDGER_URL || 'http://localhost:8080';

interface Props {
  error: string,
  setError: (err: string) => void;
  onPaymentCreated: (payment: PaymentRes) => void;
}

export const PaymentForm = ({ error = '', setError, onPaymentCreated }: Props) => {
  const inputRef = useRef<HTMLInputElement>(null);

  const [fromAccount, setFromAccount] = useState('1');
  const [toAccount, setToAccount] = useState('2');
  const [amount, setAmount] = useState('0.00');
  const [idempotencyKey, setIdempotencyKey] = useState(generateIdempotencyKey());
  const [loading, setLoading] = useState(false);

  const handleAmountChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const input = e.target;
    const cursor = input.selectionStart || 0;

    const { value, cursor: newCursor } = formatWithCursor(input.value, cursor);

    setAmount(value);

    // restore cursor AFTER render
    requestAnimationFrame(() => {
      if (inputRef.current) {
        inputRef.current.setSelectionRange(newCursor, newCursor);
      }
    });
  };

  const handleSubmit = async () => {
    setLoading(true);
    setError('');

    try {
      const res = await axios.post(`${LEDGER_URL}/payments`, {
        from_account: parseInt(fromAccount),
        to_account: parseInt(toAccount),
        amount: parseMoneyToCents(amount),
        currency: 'USD'
      }, {
        headers: { 'Idempotency-Key': idempotencyKey }
      });
      onPaymentCreated(res.data);
    } catch (er) {
      const err = er as AxiosError<{ message?: string }>;
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
          <input ref={inputRef} type="text" value={amount} onChange={handleAmountChange} className={styles.input} />
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
