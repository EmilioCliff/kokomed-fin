import { Client } from "@/data/schema";
import {
	Card,
	CardContent,
	CardDescription,
	CardFooter,
	CardHeader,
	CardTitle,
} from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { userInfo } from "os";

interface ClientCardDisplayProps {
	client: Client;
}

function ClientCardDisplay({ client }: ClientCardDisplayProps) {
	return (
		<Card className='p-2' data-user-id={client.id}>
			<CardHeader className='p-0 mb-2'>
				<CardTitle className='text-lg font-semibold'>
					{client.fullName}
				</CardTitle>
			</CardHeader>
			<div className='flex gap-x-2'>
				<div>
					<CardDescription className='text-sm text-muted-foreground'>
						Name
					</CardDescription>
					<CardContent className='p-0 text-start'>
						<p className='text-sm'>{client.fullName}</p>
					</CardContent>
				</div>
				<div>
					<CardDescription className='text-sm text-muted-foreground'>
						PhoneNumber
					</CardDescription>
					<CardContent className='p-0 text-start'>
						<p className='text-sm'>{client.phoneNumber}</p>
					</CardContent>
				</div>
			</div>
			<CardDescription className='text-sm text-muted-foreground my-2'>
				<Badge variant='outline'>{client.active ? "Active" : "Inactive"}</Badge>
			</CardDescription>
			{/* <CardContent className='p-0 text-start'>
						<p className='text-sm'>{client.fullName}</p>
					</CardContent> */}
		</Card>
	);
}

export default ClientCardDisplay;
