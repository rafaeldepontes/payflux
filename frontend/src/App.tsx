import { useState } from 'react';
import { DashboardHeader } from './components/DashboardHeader';
import { PaymentForm } from './components/PaymentForm';
import { TransactionChecker } from './components/TransactionChecker';
import { AccountBalance } from './components/AccountBalance';
import { TransactionFooter } from './components/TransactionFooter';
import type { PaymentRes } from './types';

function App() {
  const [lastPayment, setLastPayment] = useState<PaymentRes | null>(null);

  const handlePaymentCreated = (payment: PaymentRes) => {
    setLastPayment(payment);
  };

  return (
    <div className="min-h-screen bg-slate-50 text-slate-900 p-8">
      <div className="max-w-6xl mx-auto">
        <DashboardHeader />

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
          <PaymentForm onPaymentCreated={handlePaymentCreated} />
          <TransactionChecker initialId={lastPayment?.payment_id} />
          <AccountBalance lastPayment={lastPayment} />
        </div>

        <TransactionFooter payment={lastPayment} />
      </div>
    </div>
  );
}

export default App;
