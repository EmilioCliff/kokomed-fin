import { User } from '../PAGES/users/schema';
import {
	Card,
	CardContent,
	CardDescription,
	CardHeader,
	CardTitle,
} from '@/components/ui/card';

interface UserCardDisplayProps {
	user: User;
}

function UserCardDisplay({ user }: UserCardDisplayProps) {
	return (
		<Card className="p-2" data-user-id={user.id}>
			<CardHeader className="p-0 mb-2">
				<CardTitle className="text-lg font-semibold">
					{user.fullName}
				</CardTitle>
			</CardHeader>
			<div className="flex gap-x-2 justify-between">
				<div>
					<CardDescription className="text-sm text-muted-foreground">
						Phone Number
					</CardDescription>
					<CardContent className="p-0 text-start">
						<p className="text-sm">{user.phoneNumber}</p>
					</CardContent>
				</div>
				<div>
					<CardDescription className="text-sm text-muted-foreground">
						Email
					</CardDescription>
					<CardContent className="p-0 text-start">
						<p className="text-sm">{user.email}</p>
					</CardContent>
				</div>
			</div>
		</Card>
	);
}

export default UserCardDisplay;
