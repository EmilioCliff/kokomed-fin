import { type LucideIcon } from "lucide-react";
import {
	Card,
	CardContent,
	CardDescription,
	CardHeader,
	CardTitle,
} from "@/components/ui/card";

interface WidgetProps {
	title: string;
	icon: LucideIcon;
	mainAmount: number;
	active?: number;
	activeTitle: string;
	closed?: number;
	closedTitle: string;
	currency?: string;
}

function Widgets({
	title,
	icon: Icon,
	mainAmount,
	active,
	activeTitle,
	closed,
	closedTitle,
	currency,
}: WidgetProps) {
	return (
		<Card>
			<CardHeader className='flex flex-row items-center justify-between space-y-0 pb-2'>
				<CardTitle className='text-sm font-medium'>{title}</CardTitle>
				<Icon className='h-4 w-4 text-muted-foreground' />
			</CardHeader>
			<CardContent>
				<div className='text-2xl font-bold text-start'>
					{currency ? currency : ""} {mainAmount.toLocaleString()}
				</div>
				{active !== undefined && closed !== undefined && (
					<CardDescription className='text-xs flex flex-row justify-between text-muted-foreground'>
						<span>
							{activeTitle} {currency ? currency : ""} {active.toLocaleString()}
						</span>
						<span className='ml-auto'>
							{closedTitle} {currency ? currency : ""} {closed.toLocaleString()}
						</span>
					</CardDescription>
				)}
			</CardContent>
		</Card>
	);
}

export default Widgets;
