import {
	Card,
	CardHeader,
	CardTitle,
	CardDescription,
	CardContent,
} from "@/components/ui/card";
import { useState } from "react";
import { Eye, EyeOff } from "lucide-react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";

export default function Component() {
	const [showPassword, setShowPassword] = useState(false);

	return (
		<div className='grid grid-cols-2 min-h-full'>
			<div className='bg-slate-500 bg-custom-bg bg-cover bg-center'></div>
			<div className='w-full min-h-full flex justify-center items-center'>
				<Card className='mx-auto max-w-4xl border-none shadow-none flex-col text-center'>
					<CardHeader className='space-y-1 gap-4'>
						<CardTitle className='text-2xl font-bold'>LOGIN</CardTitle>
						<CardDescription className=''>
							Welcome to Kokomed Finance System
						</CardDescription>
					</CardHeader>
					<CardContent>
						<div className='space-y-4'>
							<div className='space-y-2'>
								<Input
									id='email'
									type='email'
									placeholder='m@example.com'
									required
								/>
							</div>
							<div className='space-y-2 relative'>
								<Input
									id='password'
									className='pr-10'
									placeholder='Password'
									type={showPassword ? "text" : "password"}
									required
									autoComplete='off'
								/>
								<button
									type='button'
									className='absolute right-3 top-1/3 transform -translate-y-1/2 text-gray-500 hover:text-gray-700 focus:outline-none'
									onClick={() => setShowPassword((prevState) => !prevState)}
									aria-label={showPassword ? "Hide password" : "Show password"}
								>
									{showPassword ? (
										<EyeOff className='h-5 w-5' />
									) : (
										<Eye className='h-5 w-5' />
									)}
								</button>
							</div>
							<Button type='submit' className='w-full'>
								Login
							</Button>
						</div>
					</CardContent>
				</Card>
			</div>
		</div>
	);
}
