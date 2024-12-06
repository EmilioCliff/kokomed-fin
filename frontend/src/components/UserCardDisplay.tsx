import { User } from "@/data/schema";
import {
	Card,
	CardContent,
	CardDescription,
	CardFooter,
	CardHeader,
	CardTitle,
} from "@/components/ui/card";
import { cn } from "./lib/utils";

interface UserCardDisplayProps {
	user: User;
}

function UserCardDisplay({ user }: UserCardDisplayProps) {
	return (
		<Card className='p-2' data-user-id={user.id}>
			<CardHeader className='p-0 mb-2'>
				<CardTitle className='text-lg font-semibold'>{user.fullName}</CardTitle>
			</CardHeader>
			<div className='flex gap-x-2'>
				<div>
					<CardDescription className='text-sm text-muted-foreground'>
						Name
					</CardDescription>
					<CardContent className='p-0 text-start'>
						<p className='text-sm'>{user.fullName}</p>
					</CardContent>
				</div>
				<div>
					<CardDescription className='text-sm text-muted-foreground'>
						Email
					</CardDescription>
					<CardContent className='p-0 text-start'>
						<p className='text-sm'>{user.email}</p>
					</CardContent>
				</div>
			</div>
		</Card>
	);
}

export default UserCardDisplay;
