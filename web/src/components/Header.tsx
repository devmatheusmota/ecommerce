"use client";

import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import { getToken } from "@/lib/auth";
import { useEffect, useState } from "react";
import { Search, ShoppingCart, User, LogOut, Package, Menu } from "lucide-react";

export function Header() {
  const pathname = usePathname();
  const router = useRouter();
  const [token, setToken] = useState<string | null>(null);
  const [searchQuery, setSearchQuery] = useState("");

  useEffect(() => {
    setToken(getToken());
  }, [pathname]);

  const handleLogout = () => {
    localStorage.removeItem("ecommerce_token");
    setToken(null);
    window.location.href = "/";
  };

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    if (searchQuery.trim()) {
      // For now, we just redirect to products. In a real app, we'd pass ?q=searchQuery
      router.push(`/produtos`);
    }
  };

  return (
    <header className="sticky top-0 z-50 bg-ml-yellow shadow-sm">
      <div className="mx-auto max-w-7xl px-4 py-3 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between gap-4 lg:gap-8">
          {/* Logo */}
          <Link
            href="/"
            className="flex items-center gap-2 text-2xl font-bold tracking-tight text-foreground"
          >
            <Package className="h-8 w-8" />
            <span className="hidden sm:inline">E-commerce</span>
          </Link>

          {/* Search Bar */}
          <form
            onSubmit={handleSearch}
            className="flex flex-1 max-w-2xl items-center rounded-sm bg-white shadow-sm"
          >
            <input
              type="text"
              placeholder="Buscar produtos, marcas e muito mais..."
              className="w-full px-4 py-2.5 text-sm outline-none placeholder:text-gray-400"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
            />
            <button
              type="submit"
              className="flex h-full items-center justify-center border-l border-gray-200 px-4 text-gray-500 hover:text-ml-blue"
            >
              <Search className="h-5 w-5" />
            </button>
          </form>

          {/* Right Actions */}
          <nav className="flex items-center gap-6">
            {token ? (
              <div className="hidden items-center gap-4 md:flex">
                <Link
                  href="/perfil"
                  className="flex items-center gap-2 text-sm font-medium text-foreground transition-colors hover:text-ml-blue"
                >
                  <User className="h-5 w-5" />
                  Meu Perfil
                </Link>
                <button
                  onClick={handleLogout}
                  className="flex items-center gap-2 text-sm font-medium text-foreground transition-colors hover:text-ml-blue"
                >
                  <LogOut className="h-5 w-5" />
                  Sair
                </button>
              </div>
            ) : (
              <div className="hidden items-center gap-4 md:flex">
                <Link
                  href="/registro"
                  className="text-sm font-medium text-foreground transition-colors hover:text-ml-blue"
                >
                  Crie sua conta
                </Link>
                <Link
                  href="/login"
                  className="text-sm font-medium text-foreground transition-colors hover:text-ml-blue"
                >
                  Entre
                </Link>
              </div>
            )}
            
            <Link
              href="/carrinho"
              className="relative flex items-center text-foreground transition-colors hover:text-ml-blue"
            >
              <ShoppingCart className="h-6 w-6" />
              <span className="absolute -right-2 -top-2 flex h-4 w-4 items-center justify-center rounded-full bg-ml-blue text-[10px] font-bold text-white">
                0
              </span>
            </Link>
          </nav>
        </div>

        {/* Secondary Nav */}
        <div className="mt-3 hidden items-center gap-6 text-sm font-medium text-foreground/80 md:flex">
          <Link href="/categorias" className="flex items-center gap-1 hover:text-ml-blue">
            <Menu className="h-4 w-4" />
            Categorias
          </Link>
          <Link href="/produtos" className="hover:text-ml-blue">Ofertas do dia</Link>
          <Link href="/produtos" className="hover:text-ml-blue">Histórico</Link>
          <Link href="/produtos/novo" className="hover:text-ml-blue">Vender</Link>
        </div>
      </div>
    </header>
  );
}
