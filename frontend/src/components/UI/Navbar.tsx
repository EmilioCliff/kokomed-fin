import {
	Calendar,
	ChevronDown,
	Pencil,
	User,
	Settings,
	LogOut,
	Globe,
	HelpCircle,
} from 'lucide-react';
import { Link } from 'react-router';
import { useSidebar } from '@/components/ui/sidebar';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuSeparator,
	DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { useAuth } from '@/hooks/useAuth';

function Navbar() {
	const { toggleSidebar, open } = useSidebar();
	const { logout } = useAuth();

	return (
		<>
			<header className="shared-header-sidebar-styles shadow-lg fixed top-0 left-0 z-[20] w-full flex items-center pe-2">
				<div>
					<Link
						style={{ color: 'white' }}
						className={`flex-none flex border-r-2 py-6 px-2 gap-2 ${
							open ? 'w-64' : ''
						} hover:cursor-pointer border-indigo-200/20`}
						to="/"
					>
						<Calendar />
						<h3 className="text-lg font-extrabold tracking-widest">
							Kokomed
						</h3>
					</Link>
				</div>
				<div className="border-r-2 py-6 px-4 border-indigo-200/20">
					<Calendar
						className="hover:cursor-pointer"
						onClick={toggleSidebar}
					/>
				</div>

				<p className="">Header 2</p>

				<div className="ml-auto mr-4">
					<DropdownMenu>
						<DropdownMenuTrigger
							asChild
							className="bg-transparent border-none p-0 hover:bg-transparent cursor-default"
						>
							<div className="flex flex-row border-l-2 py-4 pl-2 border-grey gap-2 align-center hover:cursor-pointer">
								<Avatar>
									<AvatarImage src="https://github.com/shadcn.png" />
									<AvatarFallback>CN</AvatarFallback>
								</Avatar>
								<p className="my-auto">JOAN</p>
								<ChevronDown size={20} className="my-auto" />
							</div>
						</DropdownMenuTrigger>
						<DropdownMenuContent>
							<DropdownMenuItem>
								<Pencil /> Edit Profile
							</DropdownMenuItem>
							<DropdownMenuItem>
								<User /> View Profile
							</DropdownMenuItem>
							<DropdownMenuSeparator />
							<DropdownMenuItem>
								<HelpCircle /> Help
							</DropdownMenuItem>
							<DropdownMenuItem>
								<Globe /> Forum
							</DropdownMenuItem>
							<DropdownMenuItem>
								<Settings /> Settings
							</DropdownMenuItem>
							<DropdownMenuSeparator />
							<DropdownMenuItem onClick={logout} className="">
								<LogOut /> Logout
							</DropdownMenuItem>
						</DropdownMenuContent>
					</DropdownMenu>
				</div>
			</header>
		</>
	);
}

export default Navbar;
