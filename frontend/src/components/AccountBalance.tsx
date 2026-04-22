import { useState } from 'react';
import axios, { AxiosError } from 'axios';
import { Wallet, CheckCircle2, AlertCircle } from 'lucide-react';
import type { BalanceRes, PaymentRes } from '../types';
import styles from './AccountBalance.module.css';

const LEDGER_URL = import.meta.env.VITE_LEDGER_URL || 'http://localhost:8080';
const RECON_URL = import.meta.env.VITE_RECON_URL || 'http://localhost:8081';

interface Props {
  setError: (err: string) => void;
  lastPayment: PaymentRes | null;
}

export const AccountBalance = ({ setError, lastPayment }: Props) => {
  const [balanceAccountId, setBalanceAccountId] = useState('1');
  const [balance, setBalance] = useState<BalanceRes | null>(null);
  const [loading, setLoading] = useState(false);

  const checkBalance = async () => {
    setError('');
    try {
      const res = await axios.get(`${LEDGER_URL}/accounts/${balanceAccountId}/balance`);
      setBalance(res.data);
    } catch (err) {
      const error = err as AxiosError<{ message?: string }>;
      setError(error.response?.data?.message || 'Failed to create payment');
    }
  };

  const createSettlement = async () => {
    if (!lastPayment) return;
    setLoading(true);
    setError('');
    try {
      await axios.post(`${RECON_URL}/settlements`, {
        transaction_id: lastPayment.payment_id,
        amount: lastPayment.amount,
        status: 'Settled'
      });
      alert('Settlement record created!');
    } catch (err) {
      const error = err as AxiosError<{ message?: string }>;
      setError(error.response?.data?.message || 'Failed to create payment');
    } finally {
      setLoading(false);
    }
  };

  const isNegative = balance ? balance.balance < 0 : false;

  return (
    <div className={styles.card}>
      <h2 className={styles.title}>
        <Wallet className={`w-5 h-5 ${isNegative ? 'text-red-500' : 'text-emerald-500'}`} />
        Account Balance
      </h2>

      <div className={styles.inputGroup}>
        <input
          type="number"
          value={balanceAccountId}
          onChange={e => setBalanceAccountId(e.target.value)}
          placeholder="Account ID"
          className={styles.input}
        />
        <button onClick={checkBalance} className={styles.checkButton}>
          Check
        </button>
      </div>

      {balance && (
        <div className={isNegative ? styles.balanceBoxNegative : styles.balanceBoxPositive}>
          <p className={isNegative ? styles.balanceLabelNegative : styles.balanceLabelPositive}>
            {isNegative ? 'Account Overdrawn' : 'Current Balance'}
          </p>
          <p className={isNegative ? styles.amountNegative : styles.amountPositive}>
            {isNegative ? `- $${(Math.abs(balance.balance) / 100).toFixed(2)}` : `$${(balance.balance / 100).toFixed(2)}`}
          </p>
          {isNegative && (
            <div className="flex items-center justify-center gap-1 mt-2 text-red-600 text-xs font-bold">
              <AlertCircle size={12} />
              Action Required
            </div>
          )}
        </div>
      )}

      <hr className={styles.divider} />

      <div className="mt-4">
        <h3 className={styles.helperTitle}>Quick Helpers</h3>
        <button
          onClick={createSettlement}
          disabled={!lastPayment || loading}
          className={styles.simulateButton}
        >
          <CheckCircle2 className="w-4 h-4" />
          Simulate Settlement
        </button>
        <p className={styles.helperText}>
          Creates an external settlement record for the last payment to enable reconciliation.
        </p>
      </div>
    </div>
  );
};