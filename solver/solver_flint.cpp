#include <solver_flint.h>
#include <stdio.h>
#include <stdlib.h>
#include <iostream>

#include <vector>
#include <algorithm>
#include <cstring>

#include <flint/flint.h>
#include <flint/fmpz.h>
#include <flint/fmpz_mod_polyxx.h>

using namespace std;
using namespace flint;

#define RET_INVALID 1
#define RET_INTERNAL_ERROR 100
#define RET_INPUT_ERROR 101

int solve_impl(vector<fmpzxx> &messages, const fmpzxx &p,  const vector<fmpzxx> &sums)
{
  vector<fmpzxx>::size_type n = sums.size();
#ifdef DEBUG
  cout << "[SOLVE-IMPL] [RECEIVED] p = " << p << endl;

  for (size_t i = 0; i < n; i++)
  {
    cout << "[SOLVE-IMPL] [RECEIVED] messages[" << i << "] = " << messages[i] << endl;
  }

  for (size_t i = 0; i < n; i++)
  {
    cout << "[SOLVE-IMPL] [RECEIVED] sums[" << i << "] = " << sums[i] << endl;
  }
#endif

  if (n < 2)
  {
#ifdef DEBUG
    cout << "Input vector too short." << endl;
#endif
    return RET_INPUT_ERROR;
  }

  // Basic sanity check to avoid weird inputs
  if (n > 1000)
  {
#ifdef DEBUG
    cout << "You probably do not want an input vector of more than 1000 elements. " << endl;
#endif
    return RET_INPUT_ERROR;
  }

  if (messages.size() != sums.size())
  {
#ifdef DEBUG
    cout << "Output vector has wrong size." << endl;
#endif
    return RET_INPUT_ERROR;
  }
#ifdef DEBUG
  cout << endl
       << "[SOLVE-IMPL] [RECEIVED] p = " << p << endl;
  cout << "[SOLVE-IMPL] [RECEIVED] n = " << n << endl
       << endl;
#endif

  if (p <= n)
  {
#ifdef DEBUG
    cout << "Prime must be (way) larger than the size of the input vector." << endl;
#endif
    return RET_INPUT_ERROR;
  }

  fmpz_mod_polyxx poly(p);
  fmpz_mod_poly_factorxx factors;
  factors.fit_length(n);
  vector<fmpzxx> coeff(n);

  // Set lead coefficient
  poly.set_coeff(n, 1);

  fmpzxx inv;
  // Compute other coeffients
  for (vector<fmpzxx>::size_type i = 0; i < n; i++)
  {
    coeff[i] = sums[i];

    vector<fmpzxx>::size_type k = 0;
    // for j = i-1, ..., 0
    for (vector<fmpzxx>::size_type j = i; j-- > 0;)
    {
      coeff[i] += coeff[k] * sums[j];
      k++;
    }
    inv = i;
    inv = -(inv + 1u);
    inv = inv.invmod(p);
    coeff[i] *= inv;
    poly.set_coeff(n - i - 1, coeff[i]);
#ifdef DEBUG
    cout << "coff[" << i << "] = " << coeff[i] << endl;
#endif
  }

#ifdef DEBUG
  cout << "Polynomial: " << endl;
  print(poly);
  cout << endl
       << endl;
#endif

  // Factor
  factors.set_factor_kaltofen_shoup(poly);

#ifdef DEBUG
  cout << "Factors: " << endl;
  print(factors);
  cout << endl
       << endl;
#endif

  vector<fmpzxx>::size_type n_roots = 0;
  for (int i = 0; i < factors.size(); i++)
  {
    if (factors.p(i).degree() != 1 || factors.p(i).lead() != 1)
    {
#ifdef DEBUG
      cout << "Non-monic factor." << endl;
#endif
      return RET_INVALID;
    }
    n_roots += factors.exp(i);
  }
  if (n_roots != n)
  {
#ifdef DEBUG
    cout << "Not enough roots." << endl;
#endif
    return RET_INVALID;
  }

  // Extract roots
  int k = 0;
  for (int i = 0; i < factors.size(); i++)
  {
    for (int j = 0; j < factors.exp(i); j++)
    {
      messages[k] = factors.p(i).get_coeff(0).negmod(p);
      k++;
    }
  }

  sort(messages.begin(), messages.end());

  return 0;
}

 char ** solve(int n, char **const out_messages, const char *prime, const char **const sums)
{
  try
  {
    fmpzxx p;

    vector<fmpzxx> s(n);
    vector<fmpzxx> messages(n);
#ifdef DEBUG
    cout << "[SOLVE] [RECEIVED] prime = " << prime << endl;
#endif

    // operator= is hard-coded to base 10 and does not check for errors
    if (fmpz_set_str(p._fmpz(), prime, 10))
    {
      return NULL;
    }


    for (size_t i = 0; i < n; i++)
    {
#ifdef DEBUG
      cout << "[SOLVE] [RECEIVED] sums[" << i << "] = " << sums[i] << endl;
#endif

      if (fmpz_set_str(s[i]._fmpz(), sums[i], 10))
      {
        return NULL;
      }

#ifdef DEBUG
      cout << "[SOLVE] [GENERATED] sums[" << i << "] = " << s[i] << endl;
#endif
    }

    for (size_t i = 0; i < n; i++)
    {
      if (out_messages[i] == NULL)
      {
        return NULL;
      }
    }
#ifdef DEBUG
    cout << "[SOLVE] [GENERATED] p = " << p << endl;
#endif

    int ret = solve_impl(messages, p, s);

    if (ret == 0)
    {
      for (size_t i = 0; i < n; i++)
      {
        // Impossible
        if (messages[i].sizeinbase(10) > strlen(prime))
        {
          return NULL;
        }
        fmpz_get_str(out_messages[i], 10, messages[i]._fmpz());
      }
    }
#ifdef DEBUG
    cout << endl;
    for (size_t i = 0; i < n; i++)
    {
      cout << "[SOLVE-IMPL] [GENERATED] messages[" << i << "] = " << messages[i] << endl;
    }
#endif

    return out_messages;
  }
  catch (...)
  {
    return NULL;
  }
}

void *allocArgv(int argc)
{
  return malloc(sizeof(char *) * argc);
}

void printArgs(int argc, char **argv)
{
  int i;
  for (i = 0; i < argc; i++)
  {
    printf("%s\n", argv[i]);
  }
}
