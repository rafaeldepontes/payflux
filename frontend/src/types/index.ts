export interface PaymentRes {
  payment_id: string;
  status: string;
  amount: number;
  currency: string;
}

export interface ReconciliationRes {
  transaction_id: string;
  status: string;
  ledger_amount: number;
  settlement_amount: number;
}

export interface RiskRes {
  transaction_id: string;
  risk_score: number;
  flags: string[];
}

export interface BalanceRes {
  account_id: number;
  balance: number;
}
