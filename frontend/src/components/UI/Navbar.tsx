import {
	ChevronDown,
	Pencil,
	User,
	Settings,
	LogOut,
	Menu,
	ChevronLeft,
	EyeOff,
	Eye,
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
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
} from '@/components/ui/dialog';
import { useAuth } from '@/hooks/useAuth';
import { useState, useEffect } from 'react';
import ThemeToggle from './ThemeToogle';
import updateUserCredentials from '@/services/updateUserCredentials';
import { toast } from 'react-toastify';
import { Input } from '../ui/input';
import { Button } from '../ui/button';

function Navbar() {
	const [formOpen, setFormOpen] = useState(false);
	const { toggleSidebar, open } = useSidebar();
	const [newPassword, setNewPassword] = useState<string>('');
	const [showPassword, setShowPassword] = useState(false);
	const { logout, accessToken, decoded } = useAuth();
	const [darkMode, setDarkMode] = useState(() => {
		return localStorage.getItem('theme') === 'dark';
	});

	const onSubmit = async (event: any) => {
		event.preventDefault();
		const response = await updateUserCredentials(newPassword, accessToken!);
		if (response.status === 'OK') {
			toast.success('Password Changed Successfully');
		} else {
			toast.error('Failed updating password');
		}

		setFormOpen(false);
		setNewPassword('');
	};

	useEffect(() => {
		if (darkMode) {
			document.documentElement.classList.add('dark');
			localStorage.setItem('theme', 'dark');
		} else {
			document.documentElement.classList.remove('dark');
			localStorage.setItem('theme', 'light');
		}
	}, [darkMode]);

	return (
		<>
			<header className="shared-header-sidebar-styles shadow-lg fixed top-0 left-0 z-[20] w-full flex items-center pe-2">
				<div>
					<Link
						style={{ color: 'white' }}
						className={`flex-none flex border-r-2 py-2 px-2 gap-2 ${
							open ? 'w-64' : ''
						} hover:cursor-pointer border-indigo-200/20`}
						to="/"
					>
						<div className="flex justify-center items-center">
							<img
								className="w-16 h-auto object-contain"
								src="/afya_credit.png"
								alt=""
							/>
							<h3 className="text-lg font-extrabold tracking-widest">
								Afya Credit
							</h3>
						</div>
					</Link>
				</div>
				<div className="border-r-2 py-6 px-4 border-indigo-200/20">
					{open ? (
						<ChevronLeft
							className="hover:cursor-pointer"
							onClick={toggleSidebar}
						/>
					) : (
						<Menu
							className="hover:cursor-pointer"
							onClick={toggleSidebar}
						/>
					)}
				</div>
				<div className="ml-auto flex items-center">
					<ThemeToggle />

					<div className="mr-4">
						<DropdownMenu>
							<DropdownMenuTrigger
								asChild
								className="bg-transparent border-none p-0 hover:bg-transparent cursor-default "
							>
								<div className="flex flex-row border-l-2 py-4 pl-2 border-grey gap-2 align-center hover:cursor-pointer">
									<Avatar>
										<AvatarImage src="https://github.com/shadcn.png" />
										<AvatarFallback>CN</AvatarFallback>
									</Avatar>
									<p className="my-auto">
										{decoded?.email
											?.split('@')[0]
											?.split('.')[0]
											.toUpperCase()}
									</p>
									<ChevronDown
										size={20}
										className="my-auto"
									/>
								</div>
							</DropdownMenuTrigger>
							<DropdownMenuContent>
								<DropdownMenuItem>
									<Pencil /> Edit Profile
								</DropdownMenuItem>
								<DropdownMenuItem
									onClick={() => setFormOpen(true)}
								>
									<EyeOff /> Change Password
								</DropdownMenuItem>
								<DropdownMenuItem>
									<User /> View Profile
								</DropdownMenuItem>
								<DropdownMenuSeparator />
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
					<Dialog open={formOpen} onOpenChange={setFormOpen}>
						<DialogContent
							aria-describedby={undefined}
							className=""
						>
							<DialogHeader>
								<DialogTitle>Change Password</DialogTitle>
							</DialogHeader>
							<form onSubmit={onSubmit}>
								<div className="space-y-2 relative">
									<Input
										placeholder="Password"
										type={
											showPassword ? 'text' : 'password'
										}
										value={newPassword}
										onChange={(e) =>
											setNewPassword(e.target.value)
										}
									/>
									<button
										type="button"
										className="absolute right-3 top-3 transform -translate-y-1/2 text-gray-500 hover:text-gray-700 focus:outline-none"
										onClick={() =>
											setShowPassword(
												(prevState) => !prevState,
											)
										}
										aria-label={
											showPassword
												? 'Hide password'
												: 'Show password'
										}
									>
										{showPassword ? (
											<EyeOff className="h-5 w-5" />
										) : (
											<Eye className="h-5 w-5" />
										)}
									</button>
									<Button
										className=" ml-auto mr-auto"
										type="submit"
									>
										Submit
									</Button>
								</div>
							</form>
							{/* <LoanForm onFormOpen={setFormOpen} /> */}
						</DialogContent>
					</Dialog>
				</div>
			</header>
		</>
	);
}

export default Navbar;
