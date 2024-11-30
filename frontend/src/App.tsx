// import "./App.css";
// import { Button } from "./components/ui/button";
// import { BrowserRouter as Router, Route, Routes } from "react-router-dom";
// import { SidebarProvider } from "@/components/ui/sidebar";
// import AppSidebar from "./components/AppSidebar";
// import Navbar from "./components/Navbar";
// import Dashbord from "./pages/Dashbord";

// function App() {
// 	return (
// 		<Router>
// 			<SidebarProvider>
// 				<AppSidebar />
// 				<Navbar />
// 				<Routes>
// 					<Route path='/' element={<Dashbord />} />
// 					<Route path='*' element={<h1 className='mt-10'>404</h1>} />
// 				</Routes>
// 			</SidebarProvider>
// 		</Router>
// 	);
// }

// export default App;
import "./App.css";
import { BrowserRouter as Router, Route, Routes } from "react-router-dom";
import { SidebarProvider } from "@/components/ui/sidebar";
import AppSidebar from "./components/AppSidebar";
import Navbar from "./components/Navbar";
import Dashboard from "./pages/Dashboard";

function App() {
	return (
		<Router>
			<SidebarProvider>
				<div className='flex'>
					<AppSidebar />
					<div className='flex-1'>
						<Navbar />
						<Routes>
							<Route path='/' element={<Dashboard />} />
							<Route path='*' element={<h1 className='mt-10'>404</h1>} />
						</Routes>
					</div>
				</div>
			</SidebarProvider>
		</Router>
	);
}

export default App;
