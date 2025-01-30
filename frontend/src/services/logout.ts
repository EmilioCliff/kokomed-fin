import api from "@/API/api";

const logoutService = async () => {
  console.log("Logging out");
  await api.get("/logout");
};

export default logoutService;
