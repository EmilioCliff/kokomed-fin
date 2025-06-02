import {
	Card,
	CardContent,
	CardDescription,
	CardHeader,
	CardTitle,
} from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { User, MapPin, Phone } from 'lucide-react';
import { formatCurrency } from '@/lib/utils';

interface ClientDetailsProps {
	client: any;
}

export function ClientDetailsCard({ client }: ClientDetailsProps) {
	return (
		<Card>
			<CardHeader>
				<CardTitle className="flex items-center gap-2">
					<User className="h-5 w-5" />
					Client Details
				</CardTitle>
				<CardDescription>
					Information about the borrower
				</CardDescription>
			</CardHeader>
			<CardContent className="space-y-4">
				<div className="flex justify-between items-start">
					<div>
						<h3 className="text-lg font-semibold">
							{client.fullName}
						</h3>
						<div className="flex items-center gap-1 text-muted-foreground">
							<Phone className="h-3.5 w-3.5" />
							<span>{client.phoneNumber}</span>
						</div>
					</div>
					<Badge variant={client.active ? 'default' : 'destructive'}>
						{client.active ? 'Active' : 'Inactive'}
					</Badge>
				</div>

				<div className="flex justify-between">
					<div>
						<p className="text-sm text-muted-foreground">Branch</p>
						<p className="text-base font-medium flex items-center gap-1">
							<MapPin className="h-3.5 w-3.5" />
							{client.branchName}
						</p>
					</div>
					<div>
						<p className="text-sm text-muted-foreground">
							ID Number
						</p>
						<p className="text-base font-medium">
							{client.idNumber || 'N/A'}
						</p>
					</div>
				</div>

				<div className="pt-2 border-t">
					<div className="flex justify-between items-center">
						<div>
							<p className="text-sm text-muted-foreground">
								Overpayment
							</p>
							<p className="text-base font-medium">
								{formatCurrency(client.overpayment)}
							</p>
						</div>
						<div>
							<p className="text-sm text-muted-foreground">
								Due Amount
							</p>
							<p className="text-base font-medium">
								{formatCurrency(client.dueAmount)}
							</p>
						</div>
					</div>
				</div>
			</CardContent>
		</Card>
	);
}
