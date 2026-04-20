import { useState, useEffect } from 'react';
import axios, { AxiosError } from 'axios';
import { Search, FileCheck, ShieldCheck } from 'lucide-react';
import type { ReconciliationRes, RiskRes } from '../types';
import styles from './TransactionChecker.module.css';

const RECON_URL = import.meta.env.VITE_RECON_URL || 'http://localhost:8081';

interface Props {
  setError: (err: string) => void;
  initialId?: string;
}

export const TransactionChecker = ({ setError, initialId = '' }: Props) => {
  const [searchId, setSearchId] = useState(initialId);
  const [reconResult, setReconResult] = useState<ReconciliationRes | null>(null);
  const [riskResult, setRiskResult] = useState<RiskRes | null>(null);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (initialId) {
      setSearchId(initialId);
      handleCheck(initialId);
    }
  }, [initialId]);

  const handleCheck = async (id: string) => {
    if (!id) return;
    setLoading(true);
    setError('');
    try {
      const [recon, risk] = await Promise.all([
        axios.get(`${RECON_URL}/reconciliation/${id}`).catch(() => ({ data: null })),
        axios.get(`${RECON_URL}/risk/${id}`).catch(() => ({ data: null }))
      ]);
      setReconResult(recon.data);
      setRiskResult(risk.data);
    } catch (err) {
      const error = err as AxiosError<{ message?: string }>;
      setError(error.response?.data?.message || 'Failed to create payment');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className={styles.card}>
      <h2 className={styles.title}>
        <Search className="w-5 h-5 text-purple-500" />
        Check Transaction
      </h2>
      <div className={styles.inputGroup}>
        <input
          type="text"
          placeholder="Payment ID"
          value={searchId}
          onChange={e => setSearchId(e.target.value)}
          className={styles.input}
        />
        <button onClick={() => handleCheck(searchId)} disabled={loading} className={styles.goButton}>
          Go
        </button>
      </div>

      <div className={styles.statusList}>
        <div className={styles.statusItem}>
          <div className={styles.statusHeader}>
            <span className={styles.statusLabel}>
              <FileCheck className="w-4 h-4" /> Reconciliation
            </span>
            {reconResult ? (
              <span className={`${styles.badge} ${reconResult.status === 'matched' ? styles.badgeMatched : styles.badgePending}`}>
                {reconResult.status === 'matched' ? 'Matched' : 'Missing Settlement'}
              </span>
            ) : <span className="text-xs text-slate-400">No data</span>}
          </div>
          {reconResult && (
            <div className="text-sm">
              <p>Ledger: <span className="font-mono font-bold">${reconResult.ledger_amount / 100}</span></p>
              <p>Settlement: <span className="font-mono font-bold">${reconResult.settlement_amount / 100}</span></p>
            </div>
          )}
        </div>

        <div className={styles.statusItem}>
          <div className={styles.statusHeader}>
            <span className={styles.statusLabel}>
              <ShieldCheck className="w-4 h-4" /> Risk Score
            </span>
            {riskResult ? (
              <span className={`${styles.badge} ${riskResult.risk_score > 30 ? styles.badgeHighRisk : styles.badgeLowRisk}`}>
                Score: {riskResult.risk_score}
              </span>
            ) : <span className="text-xs text-slate-400">No data</span>}
          </div>
          {riskResult && riskResult.flags.length > 0 && (
            <div className="flex flex-wrap gap-1 mt-2">
              {riskResult.flags.map(f => (
                <span key={f} className="text-[10px] bg-white border border-slate-200 px-2 py-0.5 rounded text-slate-600">{f}</span>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
};
