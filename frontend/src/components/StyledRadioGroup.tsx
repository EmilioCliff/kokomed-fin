import React from "react";

interface RadioOption {
	value: string;
	label: string;
}

interface StyledRadioGroupProps {
	options: RadioOption[];
	value: string;
	name: string;
	onChange: (value: string) => void;
	className?: string;
}

const StyledRadioGroup: React.FC<StyledRadioGroupProps> = ({
	options,
	value,
	name,
	onChange,
	className = "",
}) => {
	return (
		<div className={`flex space-x-4 ${className}`}>
			{options.map((option) => (
				<label key={option.value} className='flex items-center space-x-3'>
					<input
						type='radio'
						name={name}
						value={option.value}
						checked={value === option.value}
						onChange={(e) => onChange(e.target.value)}
						className='peer hidden'
					/>
					<div className='peer-checked:bg-blue-600 peer-checked:border-blue-600 h-6 w-6 rounded-full border-2 border-gray-400 bg-white flex items-center justify-center'>
						<span className='peer-checked:block hidden h-3 w-3 rounded-full bg-white'></span>
					</div>
					<span className='text-sm font-normal peer-checked:text-blue-600'>
						{option.label}
					</span>
				</label>
			))}
		</div>
	);
};

export default StyledRadioGroup;
