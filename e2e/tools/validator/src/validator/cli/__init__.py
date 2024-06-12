# SPDX-FileCopyrightText: 2024-present Sunil Thaha <sthaha@redhat.com>
#
# SPDX-License-Identifier: APACHE-2.0

import click
from validator.__about__ import __version__
from validator.stresser import ( Remote )

from validator.prometheus import MetricsValidator

from validator.cases import Cases

from validator.config import (
    Validator, load
)

pass_config = click.make_pass_decorator(Validator)

@click.group(
    context_settings={"help_option_names": ["-h", "--help"]}, 
    invoke_without_command=False,
)
@click.version_option(version=__version__, prog_name="validator")
@click.option(
   "--config-file", "-f", default="validator.yaml",
    type=click.Path(exists=True),
)
@click.pass_context
def validator(ctx: click.Context, config_file: str):
    ctx.obj = load(config_file)


@validator.command()
@click.option(
    "--script-path", "-s", 
    default="scripts/stressor.sh", 
    type=str,
)
@pass_config
def stress(cfg: Validator, script_path: str):
    remote = Remote(cfg.remote)
    result  = remote.run_script(script_path=script_path)
    click.echo(f"start_time: {result.start_time}, end_time: {result.end_time}")
    test_cases = Cases(vm = cfg.metal.vm, prom = cfg.prometheus, query_path = cfg.query_path)
    metrics_validator = MetricsValidator(cfg.prometheus)
    test_case_result, test_case_result2 = test_cases.load_test_cases()
    click.secho("Validation results during stress test:")
    for test_case in test_case_result.test_cases:

        query = test_case.refined_query

        print(f"start_time: {result.start_time}, end_time: {result.end_time}")
        metrics_res = metrics_validator.compare_metrics(result.start_time, 
                                                        result.end_time, 
                                                        query)

        click.secho(f"Query Name: {query}", fg='bright_white')
        click.secho(f"Error List: {metrics_res.el}", fg='bright_red')
        click.secho(f"Average Error: {metrics_res.me}", fg='bright_yellow')              
        
        click.secho("---------------------------------------------------", fg="cyan")
    
    click.secho("----------------------------------------------------------------------")

    for test_case in test_case_result2:

        expected_query = test_case.expected_query
        actual_query = test_case.actual_query

        print(f"start_time: {result.start_time}, end_time: {result.end_time}")
        metrics_res = metrics_validator.compare_metrics2(result.start_time, 
                                                        result.end_time, 
                                                        expected_query,
                                                        actual_query)

        click.secho(f"Query Name: EXPECTED VS ACTUAL", fg='bright_white')
        click.secho(f"Error List: {metrics_res.ape}", fg='bright_red')
        click.secho(f"Average Error: {metrics_res.mape}", fg='bright_yellow')              
        
        click.secho("---------------------------------------------------", fg="cyan")

    

