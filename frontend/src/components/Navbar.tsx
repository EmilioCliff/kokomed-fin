import {
	Calendar,
	ChevronDown,
	Pencil,
	User,
	Settings,
	LogOut,
	Info,
	Globe,
	HelpCircle,
} from "lucide-react";
import { useSidebar } from "@/components/ui/sidebar";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuLabel,
	DropdownMenuSeparator,
	DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

function Navbar() {
	const { toggleSidebar, open } = useSidebar();

	return (
		<>
			<header className='shadow-lg bg-white fixed top-0 left-0 z-[20] w-full flex items-center pe-2'>
				<div
					className={`flex-none flex border-r-2 py-6 px-2 gap-2 ${
						open ? "w-64" : ""
					}`}
				>
					<Calendar />
					<h3 className='text-lg font-extrabold tracking-widest'>Kokomed</h3>
				</div>
				<div className='border-r-2 py-6 px-4 border-grey'>
					<Calendar className='hover:cursor-pointer' onClick={toggleSidebar} />
				</div>

				<p className=''>Header 2</p>

				<div className='ml-auto mr-4'>
					<DropdownMenu>
						<DropdownMenuTrigger
							asChild
							className='bg-transparent border-none p-0 hover:bg-transparent cursor-default'
						>
							<div className='flex flex-row border-l-2 py-4 pl-2 border-grey gap-2 align-center hover:cursor-pointer'>
								<Avatar>
									<AvatarImage src='https://github.com/shadcn.png' />
									<AvatarFallback>CN</AvatarFallback>
								</Avatar>
								<p className='my-auto'>JOAN</p>
								<ChevronDown size={20} className='my-auto' />
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
							<DropdownMenuItem className=''>
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
