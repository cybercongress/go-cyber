import unittest

import humanize

from common.calculate_significance import estimate_beta, test_ranks_significance
from common.calculate_spring_rank import calculate_spring_rank, calculate_Hamiltonion, calculate_violations, \
    calculate_min_violations, calculate_system_violated_energy
from common.generate_graph import generate_graph


class GenerateGraphTest(unittest.TestCase):

    def test_generate_graph(self):
        print("Generating network")
        origin_beta = 2.1
        A, origin_raw_rank = generate_graph(N=40, beta=origin_beta, c=1)

        print("Calculating rank")
        iterations, calculated_raw_rank = calculate_spring_rank(A)
        print(f"Iterations: {iterations}")

        print("Estimate beta")
        calculated_from_origin_beta = estimate_beta(A, origin_raw_rank)
        calculated_beta = estimate_beta(A, calculated_raw_rank)
        print(f"Calculated betas: '{calculated_from_origin_beta}' and '{calculated_beta}' vs {origin_beta}")

        print("Calculating Energy")
        calculated_from_origin_energy = calculate_Hamiltonion(A, origin_raw_rank)
        calculated_energy = calculate_Hamiltonion(A, calculated_raw_rank)
        print(f"Calculated energies: '{calculated_from_origin_energy}' and '{calculated_energy}'")

        print("Calculate violations")
        v, vp, ws = calculate_violations(A, calculated_raw_rank)
        mv, mvp = calculate_min_violations(A)
        ve, vep, H = calculate_system_violated_energy(A, calculated_raw_rank)

        print(f"Violations: {v}[{vp}%] :: min violations: {mv}[{mvp}%]. Sum Aij: {ws}")
        print(f"Violation energy: {ve}[{vep}%] :: total energy: {H}")
        print("-----------------------------------------------")

        self.assertEqual(True, True)

    def test_graph_significance(self):
        A, origin_raw_rank = generate_graph(N=50, beta=0.05, c=5)
        test_ranks_significance(A, plot_file_name='test.png')
        self.assertEqual(True, True)
