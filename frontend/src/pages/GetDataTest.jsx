import React, { useState } from 'react';

const GetDataTest = () => {
	const [searchQuery, setSearchQuery] = useState('');
	const [searchType, setSearchType] = useState('phoneNumber');
	const [isLoading, setIsLoading] = useState(false);
	const [clientData, setClientData] = useState(null);
	const [error, setError] = useState(null);

	// Mock data for preview purposes
	const mockData = {
		clientDetails: {
			id: 28,
			fullName: 'Jane Doe',
			phoneNumber: '254712345678',
			branchName: 'Central Branch',
			overpayment: 150.0,
		},
		// loanDetails: {
		// 	loanId: 103,
		// 	loanAmount: 5000.0,
		// 	repayAmount: 5800.0,
		// 	disbursedOn: '2025-01-15',
		// 	dueDate: '2025-03-15',
		// 	paidAmount: 3000.0,
		// 	feePaid: true,
		// 	installments: [
		// 		{
		// 			id: 501,
		// 			installmentNumber: 1,
		// 			amountDue: 1450.0,
		// 			remainingAmount: 0,
		// 			paid: true,
		// 			paidAt: '2025-01-22',
		// 			dueDate: '2025-01-22',
		// 		},
		// 		{
		// 			id: 502,
		// 			installmentNumber: 2,
		// 			amountDue: 1450.0,
		// 			remainingAmount: 0,
		// 			paid: true,
		// 			paidAt: '2025-01-29',
		// 			dueDate: '2025-01-29',
		// 		},
		// 		{
		// 			id: 503,
		// 			installmentNumber: 3,
		// 			amountDue: 1450.0,
		// 			remainingAmount: 0,
		// 			paid: true,
		// 			paidAt: '2025-02-05',
		// 			dueDate: '2025-02-05',
		// 		},
		// 		{
		// 			id: 504,
		// 			installmentNumber: 4,
		// 			amountDue: 1450.0,
		// 			remainingAmount: 1450.0,
		// 			paid: false,
		// 			paidAt: null,
		// 			dueDate: '2025-02-12',
		// 		},
		// 	],
		// },
		paymentDetails: {
			totalPaidAmount: 3000.0,
			payments: [
				{
					id: 1,
					transactionSource: 'MPESA',
					transactionNumber: 'QXZ12345',
					accountNumber: 'ACC001',
					phoneNumber: '254712345678',
					payingName: 'JANE DOE',
					amount: 1500.0,
					paidDate: '2025-02-05T14:22:45',
					assignedTo: 28,
					assignedBy: 'APP',
				},
				{
					id: 2,
					transactionSource: 'MPESA',
					transactionNumber: 'QXZ12346',
					accountNumber: 'ACC001',
					phoneNumber: '254712345678',
					payingName: 'JANE DOE',
					amount: 1000.0,
					paidDate: '2025-01-29T10:15:30',
					assignedTo: 28,
					assignedBy: 'APP',
				},
				{
					id: 3,
					transactionSource: 'INTERNAL',
					transactionNumber: 'INT77890',
					accountNumber: 'ACC001',
					phoneNumber: '254712345678',
					payingName: 'JANE DOE',
					amount: 500.0,
					paidDate: '2025-01-22T16:45:12',
					assignedTo: 28,
					assignedBy: 'John Smith',
				},
			],
		},
	};

	const handleSearch = (e) => {
		e.preventDefault();

		if (!searchQuery.trim()) {
			setError('Please enter a client ID or phone number');
			return;
		}

		setIsLoading(true);
		setError(null);

		setTimeout(() => {
			setClientData(mockData);
			setIsLoading(false);
		}, 800);
	};

	const formatDate = (dateString) => {
		if (!dateString) return 'N/A';
		return new Date(dateString).toLocaleDateString();
	};

	const formatDateTime = (dateTimeString) => {
		if (!dateTimeString) return 'N/A';
		const date = new Date(dateTimeString);
		return `${date.toLocaleDateString()} ${date.toLocaleTimeString()}`;
	};

	return (
		<div className="bg-gray-50 min-h-screen p-4">
			<div className="max-w-6xl mx-auto">
				{/* Header */}
				<div className="bg-white p-6 rounded-lg shadow-md mb-6">
					<h1 className="text-2xl font-bold text-gray-800 mb-6">
						Client Payment Management
					</h1>

					{/* Search Form */}
					<form
						onSubmit={handleSearch}
						className="flex flex-col sm:flex-row gap-4"
					>
						<div className="flex-1">
							<input
								type="text"
								value={searchQuery}
								onChange={(e) => setSearchQuery(e.target.value)}
								placeholder={
									searchType === 'clientId'
										? 'Enter Client ID'
										: 'Enter Phone Number'
								}
								className="w-full p-3 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
							/>
						</div>
						<div className="flex-none w-full sm:w-auto">
							<select
								value={searchType}
								onChange={(e) => setSearchType(e.target.value)}
								className="w-full p-3 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
							>
								<option value="clientId">Client ID</option>
								<option value="phoneNumber">
									Phone Number
								</option>
							</select>
						</div>
						<div className="flex-none">
							<button
								type="submit"
								className="w-full bg-blue-600 hover:bg-blue-700 text-white font-medium py-3 px-6 rounded-md"
								disabled={isLoading}
							>
								{isLoading ? 'Searching...' : 'Search'}
							</button>
						</div>
					</form>

					{error && (
						<div className="mt-4 p-3 bg-red-100 text-red-700 rounded-md">
							{error}
						</div>
					)}
				</div>

				{clientData && (
					<>
						{/* Client Information */}
						<div className="bg-white p-6 rounded-lg shadow-md mb-6">
							<h2 className="text-xl font-semibold text-gray-800 mb-4">
								Client Details
							</h2>
							<div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
								<div>
									<p className="text-sm text-gray-500">
										Client ID
									</p>
									<p className="font-medium">
										{clientData.clientDetails.id || 'N/A'}
									</p>
								</div>
								<div>
									<p className="text-sm text-gray-500">
										Name
									</p>
									<p className="font-medium">
										{clientData.clientDetails.fullName ||
											'N/A'}
									</p>
								</div>
								<div>
									<p className="text-sm text-gray-500">
										Phone Number
									</p>
									<p className="font-medium">
										{clientData.clientDetails.phoneNumber ||
											'N/A'}
									</p>
								</div>
								<div>
									<p className="text-sm text-gray-500">
										Branch
									</p>
									<p className="font-medium">
										{clientData.clientDetails.branchName ||
											'N/A'}
									</p>
								</div>
								<div>
									<p className="text-sm text-gray-500">
										Overpayment
									</p>
									<p className="font-medium text-green-600">
										{clientData.clientDetails.overpayment
											? `KES ${clientData.clientDetails.overpayment.toFixed(
													2,
											  )}`
											: 'KES 0.00'}
									</p>
								</div>
							</div>
						</div>

						{/* Loan Information */}
						{clientData.loanDetails ? (
							<div className="bg-white p-6 rounded-lg shadow-md mb-6">
								<div className="flex flex-wrap justify-between items-center mb-4">
									<h2 className="text-xl font-semibold text-gray-800">
										Active Loan
									</h2>
									<div className="mt-2 lg:mt-0">
										<span className="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-blue-100 text-blue-800">
											Loan ID:{' '}
											{clientData.loanDetails.loanId}
										</span>
									</div>
								</div>

								<div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mb-6">
									<div>
										<p className="text-sm text-gray-500">
											Loan Amount
										</p>
										<p className="font-medium">
											KES{' '}
											{clientData.loanDetails.loanAmount.toFixed(
												2,
											)}
										</p>
									</div>
									<div>
										<p className="text-sm text-gray-500">
											Repayment Amount
										</p>
										<p className="font-medium">
											KES{' '}
											{clientData.loanDetails.repayAmount.toFixed(
												2,
											)}
										</p>
									</div>
									<div>
										<p className="text-sm text-gray-500">
											Amount Paid
										</p>
										<p className="font-medium">
											KES{' '}
											{clientData.loanDetails.paidAmount.toFixed(
												2,
											)}
										</p>
									</div>
									<div>
										<p className="text-sm text-gray-500">
											Disbursed On
										</p>
										<p className="font-medium">
											{formatDate(
												clientData.loanDetails
													.disbursedOn,
											)}
										</p>
									</div>
									<div>
										<p className="text-sm text-gray-500">
											Due Date
										</p>
										<p className="font-medium">
											{formatDate(
												clientData.loanDetails.dueDate,
											)}
										</p>
									</div>
									<div>
										<p className="text-sm text-gray-500">
											Processing Fee
										</p>
										<p
											className={`font-medium ${
												clientData.loanDetails.feePaid
													? 'text-green-600'
													: 'text-red-600'
											}`}
										>
											{clientData.loanDetails.feePaid
												? 'Paid'
												: 'Not Paid'}
										</p>
									</div>
								</div>

								<h3 className="text-lg font-medium text-gray-800 mb-3">
									Installments
								</h3>
								<div className="overflow-x-auto">
									<table className="min-w-full divide-y divide-gray-200">
										<thead className="bg-gray-50">
											<tr>
												<th
													scope="col"
													className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
												>
													No.
												</th>
												<th
													scope="col"
													className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
												>
													Amount Due
												</th>
												<th
													scope="col"
													className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
												>
													Remaining
												</th>
												<th
													scope="col"
													className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
												>
													Due Date
												</th>
												<th
													scope="col"
													className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
												>
													Status
												</th>
												<th
													scope="col"
													className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
												>
													Paid Date
												</th>
											</tr>
										</thead>
										<tbody className="bg-white divide-y divide-gray-200">
											{clientData.loanDetails.installments.map(
												(installment) => (
													<tr key={installment.id}>
														<td className="px-4 py-3 whitespace-nowrap text-sm">
															{
																installment.installmentNumber
															}
														</td>
														<td className="px-4 py-3 whitespace-nowrap text-sm">
															KES{' '}
															{installment.amountDue.toFixed(
																2,
															)}
														</td>
														<td className="px-4 py-3 whitespace-nowrap text-sm">
															KES{' '}
															{installment.remainingAmount.toFixed(
																2,
															)}
														</td>
														<td className="px-4 py-3 whitespace-nowrap text-sm">
															{formatDate(
																installment.dueDate,
															)}
														</td>
														<td className="px-4 py-3 whitespace-nowrap text-sm">
															<span
																className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${
																	installment.paid
																		? 'bg-green-100 text-green-800'
																		: 'bg-yellow-100 text-yellow-800'
																}`}
															>
																{installment.paid
																	? 'Paid'
																	: 'Pending'}
															</span>
														</td>
														<td className="px-4 py-3 whitespace-nowrap text-sm">
															{formatDate(
																installment.paidAt,
															)}
														</td>
													</tr>
												),
											)}
										</tbody>
									</table>
								</div>
							</div>
						) : (
							<div className="bg-white p-6 rounded-lg shadow-md mb-6">
								<div className="text-center py-6">
									<p className="text-gray-500">
										No active loans found for this client.
									</p>
								</div>
							</div>
						)}

						{/* Payment History */}
						<div className="bg-white p-6 rounded-lg shadow-md mb-6">
							<div className="flex justify-between items-center mb-4">
								<h2 className="text-xl font-semibold text-gray-800">
									Payment History
								</h2>
								<div>
									<span className="text-lg font-semibold text-blue-600">
										Total: KES{' '}
										{clientData.paymentDetails.totalPaidAmount.toFixed(
											2,
										)}
									</span>
								</div>
							</div>

							<div className="overflow-x-auto">
								<table className="min-w-full divide-y divide-gray-200">
									<thead className="bg-gray-50">
										<tr>
											<th
												scope="col"
												className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
											>
												ID
											</th>
											<th
												scope="col"
												className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
											>
												Source
											</th>
											<th
												scope="col"
												className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
											>
												Transaction No.
											</th>
											<th
												scope="col"
												className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
											>
												Account
											</th>
											<th
												scope="col"
												className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
											>
												Phone
											</th>
											<th
												scope="col"
												className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
											>
												Payer
											</th>
											<th
												scope="col"
												className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
											>
												Amount
											</th>
											<th
												scope="col"
												className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
											>
												Date
											</th>
											<th
												scope="col"
												className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
											>
												Assigned By
											</th>
										</tr>
									</thead>
									<tbody className="bg-white divide-y divide-gray-200">
										{clientData.paymentDetails.payments.map(
											(payment) => (
												<tr key={payment.id}>
													<td className="px-4 py-3 whitespace-nowrap text-sm">
														{payment.id}
													</td>
													<td className="px-4 py-3 whitespace-nowrap text-sm">
														<span
															className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${
																payment.transactionSource ===
																'MPESA'
																	? 'bg-green-100 text-green-800'
																	: 'bg-blue-100 text-blue-800'
															}`}
														>
															{
																payment.transactionSource
															}
														</span>
													</td>
													<td className="px-4 py-3 whitespace-nowrap text-sm font-medium">
														{
															payment.transactionNumber
														}
													</td>
													<td className="px-4 py-3 whitespace-nowrap text-sm">
														{payment.accountNumber}
													</td>
													<td className="px-4 py-3 whitespace-nowrap text-sm">
														{payment.phoneNumber}
													</td>
													<td className="px-4 py-3 whitespace-nowrap text-sm">
														{payment.payingName}
													</td>
													<td className="px-4 py-3 whitespace-nowrap text-sm font-medium text-green-600">
														KES{' '}
														{payment.amount.toFixed(
															2,
														)}
													</td>
													<td className="px-4 py-3 whitespace-nowrap text-sm">
														{formatDateTime(
															payment.paidDate,
														)}
													</td>
													<td className="px-4 py-3 whitespace-nowrap text-sm">
														{payment.assignedBy}
													</td>
												</tr>
											),
										)}
									</tbody>
								</table>
							</div>

							{clientData.paymentDetails.payments.length ===
								0 && (
								<div className="text-center py-6">
									<p className="text-gray-500">
										No payment records found.
									</p>
								</div>
							)}
						</div>
					</>
				)}

				{!clientData && !isLoading && !error && (
					<div className="bg-white p-6 rounded-lg shadow-md text-center">
						<p className="text-gray-500">
							Enter a client ID or phone number to view payment
							details.
						</p>
					</div>
				)}
			</div>
		</div>
	);
};

export default GetDataTest;
