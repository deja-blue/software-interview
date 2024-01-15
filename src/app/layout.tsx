import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import FlexContainer from "./components/flex-container";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "Smart Charge App",
  description: "Deja Blue Software Interview",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body className={inter.className}>
        <main className="min-h-screen bg-cover min-w-screen bg-djb">
          <FlexContainer>{children}</FlexContainer>
        </main>
      </body>
    </html>
  );
}
