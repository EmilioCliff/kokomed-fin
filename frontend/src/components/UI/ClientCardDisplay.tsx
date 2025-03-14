import { Client } from '../PAGES/customers/schema';
import {
	Card,
	CardContent,
	CardDescription,
	CardHeader,
	CardTitle,
} from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';

interface ClientCardDisplayProps {
	client: Client;
}

function ClientCardDisplay({ client }: ClientCardDisplayProps) {
	return (
		<Card className="p-2" data-user-id={client.id}>
			<CardHeader className="p-0 mb-2">
				<CardTitle className="text-lg font-semibold">
					{client.fullName}
				</CardTitle>
			</CardHeader>
			<div className="flex gap-x-2 justify-between">
				<div>
					<CardDescription className="text-sm text-muted-foreground">
						PhoneNumber
					</CardDescription>
					<CardContent className="p-0 text-start">
						<p className="text-sm">{client.phoneNumber}</p>
					</CardContent>
				</div>
				<div>
					<CardDescription className="text-sm text-muted-foreground">
						Branch
					</CardDescription>
					<CardContent className="p-0 text-start">
						<p className="text-sm">{client.branchName}</p>
					</CardContent>
				</div>
			</div>
			<CardDescription className="text-sm text-muted-foreground my-4">
				<Badge variant="outline">
					{client.active ? 'Active' : 'Inactive'}
				</Badge>
			</CardDescription>
		</Card>
	);
}

export default ClientCardDisplay;
