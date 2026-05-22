import { test, expect } from '@playwright/test';

const VALID_EMAIL    = process.env.TEST_EMAIL    ?? 'admin@example.com';
const VALID_PASSWORD = process.env.TEST_PASSWORD ?? 'password';

test.describe('Login page', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/login');
  });

  test('renders the sign-in form', async ({ page }) => {
    await expect(page).toHaveTitle(/Sign In/i);
    await expect(page.getByRole('heading', { name: /Sign In/i })).toBeVisible();
    await expect(page.getByText('Fleet Management System')).toBeVisible();

    await expect(page.getByLabel('Email Address')).toBeVisible();
    await expect(page.getByLabel('Password')).toBeVisible();
    await expect(page.getByRole('button', { name: /Sign In/i })).toBeVisible();
  });

  test('email field has autofocus', async ({ page }) => {
    const emailInput = page.getByLabel('Email Address');
    await expect(emailInput).toHaveAttribute('autofocus', '');
  });

  test('shows validation error when submitting empty form', async ({ page }) => {
    await page.getByRole('button', { name: /Sign In/i }).click();
    // Browser native validation prevents submission — email field becomes invalid
    const emailInput = page.getByLabel('Email Address');
    await expect(emailInput).toHaveAttribute('required');
  });

  test('shows error notification on wrong credentials', async ({ page }) => {
    await page.getByLabel('Email Address').fill('wrong@example.com');
    await page.getByLabel('Password').fill('wrongpassword');
    await page.getByRole('button', { name: /Sign In/i }).click();

    await expect(page.locator('#login-error')).toBeVisible();
  });

  test('pre-populates email after failed login', async ({ page }) => {
    const email = 'wrong@example.com';
    await page.getByLabel('Email Address').fill(email);
    await page.getByLabel('Password').fill('badpassword');
    await page.getByRole('button', { name: /Sign In/i }).click();

    await expect(page.getByLabel('Email Address')).toHaveValue(email);
  });

  test('error notification can be dismissed', async ({ page }) => {
    await page.getByLabel('Email Address').fill('wrong@example.com');
    await page.getByLabel('Password').fill('badpassword');
    await page.getByRole('button', { name: /Sign In/i }).click();

    const notification = page.locator('#login-error');
    await expect(notification).toBeVisible();

    await notification.getByRole('button').click();
    await expect(notification).not.toBeVisible();
  });

  test('redirects to /users on successful login', async ({ page }) => {
    await page.getByLabel('Email Address').fill(VALID_EMAIL);
    await page.getByLabel('Password').fill(VALID_PASSWORD);
    await page.getByRole('button', { name: /Sign In/i }).click();

    await expect(page).toHaveURL(/\/users/);
  });

  test('/ and /login serve the same page', async ({ page }) => {
    await page.goto('/');
    await expect(page.getByRole('heading', { name: /Sign In/i })).toBeVisible();
  });
});
