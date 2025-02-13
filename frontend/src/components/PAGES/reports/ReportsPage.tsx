import { Download, Filter, Search } from 'lucide-react';
import { useState } from 'react';
import { Button } from '@/components/ui/button';
import {
	Card,
	CardContent,
	CardDescription,
	CardFooter,
	CardHeader,
	CardTitle,
} from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { formatDistanceToNow } from 'date-fns';
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
} from '@/components/ui/dialog';
import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from '@/components/ui/select';
import { Report, ReportType } from '@/lib/types';
import { reports } from '@/data/reports';
import { ReportForm } from './ReportForm';

export function ReportsPage() {
	const [selectedReportType, setSelectedReportType] = useState<string>('all');
	const [searchQuery, setSearchQuery] = useState('');
	const [selectedReport, setSelectedReport] = useState<Report | null>(null);

	// useQuery to get the last generated time for reports

	const filteredReports = reports.filter(
		(report) =>
			(selectedReportType === 'all' ||
				report.type === selectedReportType) &&
			report.title.toLowerCase().includes(searchQuery.toLowerCase()),
	);

	return (
		<div className="px-4 space-y-6">
			<header className="flex justify-between items-center">
				<h1 className="text-3xl font-bold">Reports</h1>
			</header>

			<div className="flex flex-col sm:flex-row gap-4">
				<div className="flex-1">
					<Input
						placeholder="Search reports..."
						value={searchQuery}
						onChange={(e) => setSearchQuery(e.target.value)}
					/>
				</div>

				<Select onValueChange={setSelectedReportType}>
					<SelectTrigger className="w-full sm:w-[180px]">
						<SelectValue placeholder="Report type" />
					</SelectTrigger>
					<SelectContent>
						<SelectItem value="all">All Reports</SelectItem>
						<SelectItem value={ReportType.LOANS}>
							Loan Reports
						</SelectItem>
						<SelectItem value={ReportType.PAYMENTS}>
							Financial Reports
						</SelectItem>
						<SelectItem value={ReportType.USERS}>
							User Reports
						</SelectItem>
						<SelectItem value={ReportType.CLIENTS}>
							Customers Reports
						</SelectItem>
						<SelectItem value={ReportType.PRODUCTS}>
							Products Reports
						</SelectItem>
						<SelectItem value={ReportType.BRANCHES}>
							Branches Reports
						</SelectItem>
					</SelectContent>
				</Select>

				<Button variant="outline">
					<Filter className="mr-2 h-4 w-4" /> Filter
				</Button>
			</div>

			<div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
				{filteredReports.map((report) => (
					<Card key={report.id}>
						<CardHeader>
							<CardTitle>{report.title}</CardTitle>
							<CardDescription>
								{report.description}
							</CardDescription>
						</CardHeader>
						<CardContent>
							{/* <p>Last generated: {timeAgo(report.lastGenerated)}</p> */}
							<p>Last generated: {report.lastGenerated}</p>
						</CardContent>
						<CardFooter>
							<Dialog
								open={selectedReport?.id === report.id}
								onOpenChange={(open) =>
									setSelectedReport(open ? report : null)
								}
							>
								<DialogTrigger asChild>
									<Button
										className="w-full"
										size="sm"
										variant="outline"
									>
										Generate Report
									</Button>
								</DialogTrigger>
								<DialogContent className="max-w-screen-lg">
									<DialogHeader>
										<DialogTitle>
											Generate {selectedReport?.title}
										</DialogTitle>
										<DialogDescription>
											Enter the required details to
											generate the report.
										</DialogDescription>
									</DialogHeader>

									{selectedReport && (
										<ReportForm
											report={selectedReport}
											onClose={() =>
												setSelectedReport(null)
											}
										/>
									)}
								</DialogContent>
							</Dialog>
						</CardFooter>
					</Card>
				))}
			</div>
		</div>
	);
}

function timeAgo(dateString: string): string {
	const date = new Date(dateString);
	return formatDistanceToNow(date, { addSuffix: true });
}
