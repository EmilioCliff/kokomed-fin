import { useParams } from 'react-router';

function PaymentLoanBreakdown() {
	const { id } = useParams();
	return <div>PaymentLoanBreakdown {id}</div>;
}

export default PaymentLoanBreakdown;
