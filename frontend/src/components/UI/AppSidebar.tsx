import {
	Calendar,
	Home,
	Inbox,
	Search,
	Settings,
	ChevronRight,
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
		icon: Home,
		links: [
			{
				title: 'Overview',
				url: '/loans/overview',
			},
			{
				title: 'Loans Timeline',
				url: '#',
			},
			{
				title: 'Reports',
				url: '#',
			},
		],
	},
	{
		title: 'Payments Module',
		url: '#',
		icon: Calendar,
		links: [
			{
				title: 'Overview',
				url: '/payments/overview',
			},
			{
				title: 'Payments Calendar',
				url: '#',
			},
			{
				title: 'Reports',
				url: '#',
			},
		],
	},
	{
		title: 'Products Module',
		url: '#',
		icon: Calendar,
		links: [
			{
				title: 'Overview',
				url: '/products/overview',
			},
			{
				title: 'Reports',
				url: '#',
			},
		],
	},
	{
		title: 'Customers Module',
		url: '#',
		icon: Inbox,
		links: [
			{
				title: 'Overview',
				url: '/customers/overview',
			},
			{
				title: 'Reports',
				url: '#',
			},
		],
	},
	{
		title: 'Users Module',
		url: '#',
		icon: Search,
		links: [
			{
				title: 'Overview',
				url: '/users/overview',
			},
			{
				title: 'Reports',
				url: '#',
			},
		],
	},
	{
		title: 'Branches',
		url: '#',
		icon: Settings,
		links: [
			{
				title: 'Overview',
				url: '/branches/overview',
			},
			{
				title: 'Reports',
				url: '#',
			},
		],
	},
];
