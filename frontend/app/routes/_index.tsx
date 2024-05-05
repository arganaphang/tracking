import type { MetaFunction } from "@remix-run/node";

export const meta: MetaFunction = () => {
  return [
    { title: "Order Tracking System" },
    { name: "description", content: "Welcome to Tracking System" },
  ];
};

export default function Index() {
  return (
    <div className="w-full min-h-screen flex flex-col justify-center items-center">
      <h1>Hello World</h1>
    </div>
  );
}
