import {
	Calendar,
	Home,
	Inbox,
	Search,
	Settings,
	ChevronRight,
} from "lucide-react";
import {
	Sidebar,
	SidebarContent,
	SidebarGroup,
	SidebarHeader,
	SidebarMenu,
	SidebarMenuItem,
	SidebarMenuBadge,
	SidebarMenuAction,
	SidebarMenuButton,
	SidebarMenuSub,
	SidebarMenuSubItem,
	SidebarMenuSubButton,
	SidebarGroupLabel,
	SidebarGroupAction,
	SidebarGroupContent,
	useSidebar,
} from "@/components/ui/sidebar";
import {
	Collapsible,
	CollapsibleTrigger,
	CollapsibleContent,
} from "./ui/collapsible";
import { Separator } from "./ui/separator";

const items = [
	{
		title: "Loans Module",
		url: "#",
		icon: Home,
		links: [
			{
				title: "Overview",
				url: "/loans/overview",
			},
			{
				title: "Loans Timeline",
				url: "#",
			},
			{
				title: "Reports",
				url: "#",
			},
		],
	},
	{
		title: "Payments Module",
		url: "#",
		icon: Calendar,
		links: [
			{
				title: "Overview",
				url: "/payments/overview",
			},
			{
				title: "Payments Calendar",
				url: "#",
			},
			{
				title: "Reports",
				url: "#",
			},
		],
	},
	{
		title: "Products Module",
		url: "#",
		icon: Calendar,
		links: [
			{
				title: "Overview",
				url: "/products/overview",
			},
			{
				title: "Reports",
				url: "#",
			},
		],
	},
	{
		title: "Customers Module",
		url: "#",
		icon: Inbox,
		links: [
			{
				title: "Overview",
				url: "/customers/overview",
			},
			{
				title: "Reports",
				url: "#",
			},
		],
	},
	{
		title: "Users Module",
		url: "#",
		icon: Search,
		links: [
			{
				title: "Overview",
				url: "/users/overview",
			},
			{
				title: "Reports",
				url: "#",
			},
		],
	},
	{
		title: "Branches",
		url: "#",
		icon: Settings,
		links: [
			{
				title: "Overview",
				url: "/branches/overview",
			},
			{
				title: "Reports",
				url: "#",
			},
		],
	},
];

function AppSidebar() {
	const { open } = useSidebar();

	return (
		<>
			<Sidebar
				className='mt-16 pt-4'
				side='left'
				variant='sidebar'
				collapsible='icon'
			>
				<SidebarHeader className='flex flex-row'>
					<Home width={24} className='mx-auto flex-none' />
					{open ? <h3 className='grow text-start'>Loan Management</h3> : ""}
				</SidebarHeader>
				<Separator className='mt-2 mb-4' />
				<SidebarContent>
					<SidebarGroup>
						<SidebarGroupContent>
							<SidebarMenu>
								{items.map((item, index) => (
									<Collapsible
										key={index}
										{...(index === 0 && { defaultOpen: true })}
										className='group/collapsible'
									>
										<SidebarMenuItem key={item.title}>
											<CollapsibleTrigger asChild>
												<SidebarMenuButton asChild className='text-black'>
													{/* isActive */}
													<a href={item.url}>
														<item.icon />
														<span>{item.title}</span>
													</a>
												</SidebarMenuButton>
											</CollapsibleTrigger>
											<CollapsibleContent>
												<SidebarMenuSub>
													{item.links.map((link, index) => (
														<SidebarMenuSubItem key={index}>
															<SidebarMenuSubButton asChild>
																<a href={link.url}>
																	<span>{link.title}</span>
																</a>
															</SidebarMenuSubButton>
														</SidebarMenuSubItem>
													))}
												</SidebarMenuSub>
											</CollapsibleContent>
											<SidebarMenuBadge>
												<ChevronRight className='ml-auto transition-transform group-data-[state=open]/collapsible:rotate-90' />
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
