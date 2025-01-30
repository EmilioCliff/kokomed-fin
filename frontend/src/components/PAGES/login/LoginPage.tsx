import {
  Card,
  CardHeader,
  CardTitle,
  CardDescription,
  CardContent,
} from "@/components/ui/card";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Eye, EyeOff } from "lucide-react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { toast } from "react-toastify";
import { useNavigate } from "react-router";
import { loginFormSchema, LoginForm } from "./schema";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormMessage,
} from "@/components/ui/form";
import { useAuth } from "@/hooks/useAuth";

export default function LoginPage() {
  const { login } = useAuth();

  const [showPassword, setShowPassword] = useState(false);
  const navigate = useNavigate();

  function onError(errors: any) {
    console.log("Errors: ", errors);
  }

  const onSubmit = async (values: LoginForm) => {
    try {
      const data = await login(values.email, values.password);
      if (data) {
        navigate("/");
      }
    } catch (error: any) {
      toast.error(error.message);
    }
  };

  const form = useForm<LoginForm>({
    resolver: zodResolver(loginFormSchema),
    defaultValues: {
      email: "",
      password: "",
    },
  });

  return (
    <>
      {/* {mutation.isPending && <Spinner />} */}
      <div className="grid grid-cols-1 min-h-full sm:grid-cols-1 md:grid-cols-2">
        <div className="bg-slate-500 bg-custom-bg bg-cover bg-center somen"></div>
        <div className="w-full min-h-full flex justify-center items-center">
          <Card className="mx-auto max-w-4xl border-none shadow-none flex-col text-center">
            <CardHeader className="space-y-1 gap-4">
              <CardTitle className="text-2xl font-bold">LOGIN</CardTitle>
              <CardDescription className="">
                Welcome to Kokomed Finance System
              </CardDescription>
            </CardHeader>
            <CardContent>
              <Form {...form}>
                <form
                  onSubmit={form.handleSubmit(onSubmit, onError)}
                  className="space-y-8"
                >
                  <FormField
                    control={form.control}
                    name="email"
                    render={({ field }) => (
                      <FormItem>
                        <FormControl>
                          <Input placeholder="Email" {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                  <FormField
                    control={form.control}
                    name="password"
                    render={({ field }) => (
                      <FormItem>
                        <FormControl>
                          <div className="space-y-2 relative">
                            <Input
                              placeholder="Password"
                              type={showPassword ? "text" : "password"}
                              {...field}
                            />
                            <button
                              type="button"
                              className="absolute right-3 top-1/3 transform -translate-y-1/2 text-gray-500 hover:text-gray-700 focus:outline-none"
                              onClick={() => setShowPassword((prevState) => !prevState)}
                              aria-label={
                                showPassword ? "Hide password" : "Show password"
                              }
                            >
                              {showPassword ? (
                                <EyeOff className="h-5 w-5" />
                              ) : (
                                <Eye className="h-5 w-5" />
                              )}
                            </button>
                          </div>
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <Button className=" ml-auto mr-auto" type="submit">
                    Submit
                  </Button>
                </form>
              </Form>
              {/* <div className="space-y-4">
                <div className="space-y-2">
                  <Input
                    id="email"
                    value={credentials.email}
                    onChange={handleChange}
                    type="email"
                    placeholder="m@example.com"
                    required
                  />
                </div>
                <div className="space-y-2 relative">
                  <Input
                    id="password"
                    className="pr-10"
                    placeholder="Password"
                    type={showPassword ? 'text' : 'password'}
                    required
                    autoComplete="off"
                    value={credentials.password}
                    onChange={handleChange}
                  />
                  <button
                    type="button"
                    className="absolute right-3 top-1/3 transform -translate-y-1/2 text-gray-500 hover:text-gray-700 focus:outline-none"
                    onClick={() => setShowPassword((prevState) => !prevState)}
                    aria-label={showPassword ? 'Hide password' : 'Show password'}
                  >
                    {showPassword ? (
                      <EyeOff className="h-5 w-5" />
                    ) : (
                      <Eye className="h-5 w-5" />
                    )}
                  </button>
                </div>
                <Button onClick={onSubmit} type="submit" className="w-full">
                  Login
                </Button>
              </div> */}
            </CardContent>
          </Card>
        </div>
      </div>
    </>
  );
}
