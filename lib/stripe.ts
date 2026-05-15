import Stripe from "stripe";

export const stripe = new Stripe(process.env.STRIPE_SECRET_KEY ?? "sk_test_placeholder");

export const priceIds = {
  BASIC: process.env.STRIPE_BASIC_PRICE_ID ?? "price_basic_monthly",
  PREMIUM: process.env.STRIPE_PREMIUM_PRICE_ID ?? "price_premium_monthly"
};
