import {
	ChevronRight,
	Landmark,
	CreditCard,
	Package,
	Users,
	User,
	Building,
	FileBarChart,
} from 'lucide-react';
import { Link, useLocation } from 'react-router';
import {
	Sidebar,
	SidebarContent,
	SidebarGroup,
	SidebarMenu,
	SidebarMenuItem,
	SidebarMenuBadge,
	SidebarMenuButton,
	SidebarMenuSub,
	SidebarMenuSubItem,
	SidebarMenuSubButton,
	SidebarGroupContent,
} from '@/components/ui/sidebar';
import {
	Collapsible,
	CollapsibleTrigger,
	CollapsibleContent,
} from '../ui/collapsible';
import { useTable } from '@/hooks/useTable';
import { title } from 'process';

function AppSidebar() {
	const { resetTableState } = useTable();

	return (
		<>
			<Sidebar
				className="mt-16 pt-4"
				side="left"
				variant="sidebar"
				collapsible="icon"
			>
				<SidebarContent>
					<SidebarGroup>
						<SidebarGroupContent>
							<SidebarMenu>
								{items.map((item, index) => (
									<Collapsible
										key={index}
										{...(index === 0 && {
											defaultOpen: true,
										})}
										className="group/collapsible"
									>
										<SidebarMenuItem key={item.title}>
											<CollapsibleTrigger asChild>
												<SidebarMenuButton
													asChild
													className="text-black"
												>
													<a
														style={{
															color: 'white',
														}}
														href={item.url}
													>
														<item.icon />
														<span>
															{item.title}
														</span>
													</a>
												</SidebarMenuButton>
											</CollapsibleTrigger>
											<CollapsibleContent>
												<SidebarMenuSub>
													{item.links.map(
														(link, index) => {
															const location =
																useLocation();
															const isActive =
																location.pathname ===
																link.url;

															return (
																<SidebarMenuSubItem
																	key={index}
																>
																	<SidebarMenuSubButton
																		className={`${
																			isActive
																				? 'bg-sidebar-accent text-sidebar-accent-foreground'
																				: 'bg-transparent text-gray-400'
																		} `}
																		asChild
																	>
																		<Link
																			onClick={
																				resetTableState
																			}
																			to={
																				link.url
																			}
																		>
																			<span>
																				{
																					link.title
																				}
																			</span>
																		</Link>
																	</SidebarMenuSubButton>
																</SidebarMenuSubItem>
															);
														},
													)}
												</SidebarMenuSub>
											</CollapsibleContent>
											<SidebarMenuBadge>
												<ChevronRight className="ml-auto transition-transform group-data-[state=open]/collapsible:rotate-90" />
											</SidebarMenuBadge>
										</SidebarMenuItem>
									</Collapsible>
								))}
							</SidebarMenu>
						</SidebarGroupContent>
					</SidebarGroup>
					<SidebarGroup />
				</SidebarContent>
			</Sidebar>
		</>
	);
}

export default AppSidebar;

const items = [
	{
		title: 'Loans Module',
		url: '#',
		icon: Landmark,
		links: [
			{
				title: 'Overview',
				url: '/loans/overview',
			},
			{
				title: 'Expected Payments',
				url: '/loans/expected-payments',
			},
			{
				title: 'Loans Timeline',
				url: '/loans/timeline',
			},
		],
	},
	{
		title: 'Payments Module',
		url: '#',
		icon: CreditCard,
		links: [
			{
				title: 'Overview',
				url: '/payments/overview',
			},
			{
				title: 'Payments Calendar',
				url: '#',
			},
		],
	},
	{
		title: 'Products Module',
		url: '#',
		icon: Package,
		links: [
			{
				title: 'Overview',
				url: '/products/overview',
			},
		],
	},
	{
		title: 'Customers Module',
		url: '#',
		icon: Users,
		links: [
			{
				title: 'Overview',
				url: '/customers/overview',
			},
		],
	},
	{
		title: 'Users Module',
		url: '#',
		icon: User,
		links: [
			{
				title: 'Overview',
				url: '/users/overview',
			},
		],
	},
	{
		title: 'Branches',
		url: '#',
		icon: Building,
		links: [
			{
				title: 'Overview',
				url: '/branches/overview',
			},
		],
	},
	{
		title: 'Reports',
		url: '#',
		icon: FileBarChart,
		links: [
			{
				title: 'Reports',
				url: '/reports',
			},
		],
	},
];
